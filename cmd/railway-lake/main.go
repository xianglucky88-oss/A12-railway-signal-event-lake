package main

import (
	"flag"
	"fmt"
	"os"

	"example.com/railway-signal-lake/internal/alarm"
	"example.com/railway-signal-lake/internal/analysis"
	"example.com/railway-signal-lake/internal/demo"
	"example.com/railway-signal-lake/internal/model"
	"example.com/railway-signal-lake/internal/service"
	"example.com/railway-signal-lake/internal/store"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
func run(args []string) error {
	if len(args) == 0 {
		usage()
		return nil
	}
	switch args[0] {
	case "ingest":
		return ingest(args[1:])
	case "compute":
		return compute(args[1:])
	case "demo":
		return runDemo(args[1:])
	case "help", "-h", "--help":
		usage()
		return nil
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}
func usage() {
	fmt.Print(`railway-lake - railway signal event lake
Commands:
  ingest  --events FILE [--links FILE] [--work-orders FILE] --data-dir DIR
  compute --data-dir DIR
  demo    [--data-dir DIR]
`)
}
func ingest(args []string) error {
	fs := flag.NewFlagSet("ingest", flag.ContinueOnError)
	events := fs.String("events", "", "event JSONL")
	links := fs.String("links", "", "topology link JSONL")
	orders := fs.String("work-orders", "", "work order JSONL")
	dir := fs.String("data-dir", "data", "lake root")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *events == "" && *links == "" && *orders == "" {
		return fmt.Errorf("at least one input file is required")
	}
	lake := service.New(*dir)
	if *events != "" {
		rows, err := store.ReadJSONL[model.Event](*events)
		if err != nil {
			return err
		}
		if err := lake.IngestEvents(rows); err != nil {
			return err
		}
	}
	if *links != "" {
		rows, err := store.ReadJSONL[model.TopologyLink](*links)
		if err != nil {
			return err
		}
		if err := lake.IngestLinks(rows); err != nil {
			return err
		}
	}
	if *orders != "" {
		rows, err := store.ReadJSONL[model.WorkOrder](*orders)
		if err != nil {
			return err
		}
		if err := lake.IngestWorkOrders(rows); err != nil {
			return err
		}
	}
	return writeJSON(map[string]string{"status": "ok"})
}
func compute(args []string) error {
	fs := flag.NewFlagSet("compute", flag.ContinueOnError)
	dir := fs.String("data-dir", "data", "lake root")
	if err := fs.Parse(args); err != nil {
		return err
	}
	result, err := service.New(*dir).Compute(analysis.Config{
		Alarm:               alarm.Config{DuplicateWindowMs: 60_000, ChatterWindowMs: 120_000, ChatterThreshold: 3, StormWindowMs: 30_000, StormDeviceCount: 2},
		CorrelationWindowMs: 60_000, WorkOrderToleranceMs: 120_000,
	})
	if err != nil {
		return err
	}
	return writeJSON(result)
}
func runDemo(args []string) error {
	fs := flag.NewFlagSet("demo", flag.ContinueOnError)
	dir := fs.String("data-dir", ".demo-data", "demo lake root")
	if err := fs.Parse(args); err != nil {
		return err
	}
	lake := service.New(*dir)
	if err := lake.IngestEvents(demo.Events()); err != nil {
		return err
	}
	if err := lake.IngestLinks(demo.Links()); err != nil {
		return err
	}
	if err := lake.IngestWorkOrders(demo.WorkOrders()); err != nil {
		return err
	}
	result, err := lake.Compute(analysis.Config{Alarm: alarm.Config{}, CorrelationWindowMs: 60_000, WorkOrderToleranceMs: 120_000})
	if err != nil {
		return err
	}
	fmt.Printf("demo complete: events=%d incidents=%d correlations=%d\n", len(result.Snapshots), len(result.Incidents), len(result.Correlations))
	return writeJSON(result)
}
