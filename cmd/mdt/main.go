package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/ecc1/medtronic"
)

type prog func(*medtronic.Pump, []string) interface{}

var command = map[string]prog{
	"basal":         basal,
	"battery":       battery,
	"carbratios":    carbRatios,
	"carbunits":     carbUnits,
	"clock":         clock,
	"execute":       execute,
	"glucoseunits":  glucoseUnits,
	"history":       history,
	"model":         model,
	"pumpid":        pumpID,
	"reservoir":     reservoir,
	"resume":        resume,
	"rssi":          rssi,
	"sensitivities": sensitivities,
	"setclock":      setClock,
	"settempbasal":  setTempBasal,
	"settings":      settings,
	"status":        status,
	"suspend":       suspend,
	"targets":       targets,
	"tempbasal":     tempBasal,
	"wakeup":        wakeup,
}

// TODO: add per-command help

func usage() {
	eprintf("Usage: %s command [arg ...]\n", os.Args[0])
	eprintf("Commands:")
	keys := make([]string, len(command))
	i := 0
	for k := range command {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	for _, k := range keys {
		eprintf(" %s", k)
	}
	eprintf("\n")
	os.Exit(1)
}

func main() {
	if len(os.Args) == 1 {
		usage()
	}
	name := os.Args[1]
	args := os.Args[2:]
	prog := command[name]
	if prog == nil {
		eprintf("%s: unknown command\n", name)
		usage()
	}
	pump := medtronic.Open()
	defer pump.Close()
	pump.Wakeup()
	result := prog(pump, args)
	if pump.Error() != nil {
		log.Fatal(pump.Error())
	}
	if result == nil {
		return
	}
	b, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Println(err)
		fmt.Println(result)
		return
	}
	fmt.Println(string(b))
}

func eprintf(format string, arg ...interface{}) {
	fmt.Fprintf(os.Stderr, format, arg...) // nolint
}

// TODO: with argument to schedule progs, get schedule at that time

func basal(pump *medtronic.Pump, _ []string) interface{} {
	return pump.BasalRates()
}

func battery(pump *medtronic.Pump, _ []string) interface{} {
	return pump.Battery()
}

func carbRatios(pump *medtronic.Pump, _ []string) interface{} {
	return pump.CarbRatios()
}

func carbUnits(pump *medtronic.Pump, _ []string) interface{} {
	return pump.CarbUnits()
}

func clock(pump *medtronic.Pump, _ []string) interface{} {
	return pump.Clock()
}

func execute(pump *medtronic.Pump, args []string) interface{} {
	if len(args) == 0 {
		executeUsage()
	}
	cmd, err := strconv.ParseUint(args[0], 16, 8)
	if err != nil {
		eprintf("%v\n", err)
		executeUsage()
	}
	var params []byte
	for _, s := range args[1:] {
		b, err := strconv.ParseUint(s, 16, 8)
		if err != nil {
			eprintf("%v\n", err)
			executeUsage()
		}
		params = append(params, byte(b))
	}
	return pump.Execute(medtronic.Command(cmd), params...)
}

func executeUsage() {
	eprintf("Usage: execute cmd [param ...]\n")
	os.Exit(1)
}

func glucoseUnits(pump *medtronic.Pump, _ []string) interface{} {
	return pump.GlucoseUnits()
}

func history(pump *medtronic.Pump, _ []string) interface{} {
	return pump.LastHistoryPage()
}

func model(pump *medtronic.Pump, _ []string) interface{} {
	return pump.Model()
}

func pumpID(pump *medtronic.Pump, _ []string) interface{} {
	return pump.PumpID()
}

func reservoir(pump *medtronic.Pump, _ []string) interface{} {
	return pump.Reservoir()
}

func resume(pump *medtronic.Pump, _ []string) interface{} {
	log.Printf("resuming pump")
	pump.Suspend(false)
	return nil
}

func rssi(pump *medtronic.Pump, _ []string) interface{} {
	return pump.RSSI()
}

func sensitivities(pump *medtronic.Pump, _ []string) interface{} {
	return pump.InsulinSensitivities()
}

func setClock(pump *medtronic.Pump, args []string) interface{} {
	t := time.Time{}
	switch len(args) {
	case 2:
		t = parseTime(args[0] + " " + args[1])
	case 1:
		if args[0] == "now" {
			t = time.Now()
		} else {
			setClockUsage()
		}
	default:
		setClockUsage()
	}
	log.Printf("setting pump clock to %s", t.Format(medtronic.UserTimeLayout))
	pump.SetClock(t)
	return nil
}

func parseTime(date string) time.Time {
	t, err := time.ParseInLocation(medtronic.UserTimeLayout, date, time.Local)
	if err != nil {
		setClockUsage()
	}
	return t
}

func setClockUsage() {
	eprintf("Usage: setclock YYYY-MM-DD HH:MM:SS (or \"now\")\n")
	os.Exit(1)
}

func setTempBasal(pump *medtronic.Pump, args []string) interface{} {
	if len(args) != 2 || len(args[1]) == 0 {
		setTempBasalUsage()
	}
	duration, err := time.ParseDuration(args[0])
	if err != nil {
		setTempBasalUsage()
	}
	rateArg := args[1]
	n := len(rateArg) - 1
	if rateArg[n] == '%' {
		percent, err := strconv.Atoi(rateArg[:n])
		if err != nil {
			setTempBasalUsage()
		}
		log.Printf("setting temporary basal of %d%% for %v", percent, duration)
		pump.SetPercentTempBasal(duration, percent)
	} else {
		f, err := strconv.ParseFloat(rateArg, 32)
		if err != nil {
			setTempBasalUsage()
		}
		rate := medtronic.Insulin(1000.0*f + 0.5)
		log.Printf("setting temporary basal of %d.%03d units/hour for %v", rate/1000, rate%1000, duration)
		pump.SetAbsoluteTempBasal(duration, rate)
	}
	return nil
}

func setTempBasalUsage() {
	eprintf("Usage: settempbasal duration (units/hr | rate%%)\n")
	os.Exit(1)
}

func settings(pump *medtronic.Pump, _ []string) interface{} {
	return pump.Settings()
}

func status(pump *medtronic.Pump, _ []string) interface{} {
	return pump.Status()
}

func suspend(pump *medtronic.Pump, _ []string) interface{} {
	log.Printf("suspending pump")
	pump.Suspend(true)
	return nil
}

func targets(pump *medtronic.Pump, _ []string) interface{} {
	return pump.GlucoseTargets()
}

func tempBasal(pump *medtronic.Pump, _ []string) interface{} {
	return pump.TempBasal()
}

func wakeup(pump *medtronic.Pump, _ []string) interface{} {
	// pump.Wakeup has already been called
	return nil
}
