package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ecc1/medtronic"
)

type (
	// Prog represents a function to be executed with the given arguments on the pump.
	Prog func(*medtronic.Pump, Arguments) interface{}

	// Command specifies the function and formal parameters for a command.
	// If Variadic is true, the last parameter is bound to a list of arguments
	// (a JSON array of strings, or the remaining command-line arguments).
	Command struct {
		Cmd      Prog
		Params   []string
		Variadic bool
	}
)

var (
	// TODO: add per-command help
	command = map[string]Command{
		"basal":         cmd(basal),
		"battery":       cmd(battery),
		"bolus":         cmd(bolus, "units"),
		"button":        cmdN(button, "keys"),
		"carbratios":    cmd(carbRatios),
		"carbunits":     cmd(carbUnits),
		"clock":         cmd(clock),
		"execute":       cmdN(execute, "command", "arguments"),
		"firmware":      cmd(firmware),
		"glucoseunits":  cmd(glucoseUnits),
		"model":         cmd(model),
		"pumpid":        cmd(pumpID),
		"reservoir":     cmd(reservoir),
		"resume":        cmd(resume),
		"rssi":          cmd(rssi),
		"sensitivities": cmd(sensitivities),
		"setclock":      cmd(setClock, "time"),
		"setmaxbolus":   cmd(setMaxBolus, "units"),
		"settempbasal":  cmd(setTempBasal, "temp", "rate", "duration"),
		"settings":      cmd(settings),
		"status":        cmd(status),
		"suspend":       cmd(suspend),
		"targets":       cmd(targets),
		"tempbasal":     cmd(tempBasal),
		"wakeup":        cmd(wakeup),
	}
)

func cmd(prog Prog, params ...string) Command {
	return Command{Cmd: prog, Params: params, Variadic: false}
}

func cmdN(prog Prog, params ...string) Command {
	return Command{Cmd: prog, Params: params, Variadic: true}
}

func cmdError(name string, msg string, err error) {
	eprintf("%s: %v\n", name, err)
	eprintf("usage: %s %s\n", name, msg)
	os.Exit(1)
}

// TODO: with argument to schedule progs, get schedule at that time

func basal(pump *medtronic.Pump, _ Arguments) interface{} {
	return pump.BasalRates()
}

func battery(pump *medtronic.Pump, _ Arguments) interface{} {
	return pump.Battery()
}

func bolus(pump *medtronic.Pump, args Arguments) interface{} {
	f, err := args.Float("units")
	if err != nil {
		cmdError("bolus", "units", err)
	}
	amount := medtronic.Insulin(1000.0*f + 0.5)
	log.Printf("performing bolus of %v units", amount)
	pump.Bolus(amount)
	return nil
}

func button(pump *medtronic.Pump, args Arguments) interface{} {
	v, err := args.Strings("keys")
	if err != nil {
		buttonUsage(err)
	}
	for _, s := range v {
		b := parseButton(s)
		log.Printf("pressing %v", b)
		pump.Button(b)
		if pump.Error() != nil {
			break
		}
	}
	return nil
}

var buttonName = map[string]medtronic.PumpButton{
	"b":    medtronic.BolusButton,
	"esc":  medtronic.EscButton,
	"act":  medtronic.ActButton,
	"up":   medtronic.UpButton,
	"down": medtronic.DownButton,
}

func parseButton(s string) medtronic.PumpButton {
	b, found := buttonName[s]
	if !found {
		buttonUsage(fmt.Errorf("unknown pump button %q", s))
	}
	return b
}

func buttonUsage(err error) {
	cmdError("button", "(b|esc|act|up|down) ...", err)
}

func carbRatios(pump *medtronic.Pump, _ Arguments) interface{} {
	return pump.CarbRatios()
}

func carbUnits(pump *medtronic.Pump, _ Arguments) interface{} {
	return pump.CarbUnits()
}

func clock(pump *medtronic.Pump, _ Arguments) interface{} {
	return pump.Clock()
}

func execute(pump *medtronic.Pump, args Arguments) interface{} {
	c, err := strconv.ParseUint(args["command"].(string), 16, 8)
	if err != nil {
		executeUsage(err)
	}
	v, err := args.Strings("arguments")
	if err != nil {
		executeUsage(err)
	}
	var params []byte
	for _, s := range v {
		b, err := strconv.ParseUint(s, 16, 8)
		if err != nil {
			executeUsage(err)
		}
		params = append(params, byte(b))
	}
	cmd := medtronic.Command(c)
	log.Printf("executing %v % X", cmd, params)
	return pump.Execute(cmd, params...)
}

func executeUsage(err error) {
	cmdError("execute", "cmd [param ...]", err)
}

func firmware(pump *medtronic.Pump, _ Arguments) interface{} {
	return pump.FirmwareVersion()
}

func glucoseUnits(pump *medtronic.Pump, _ Arguments) interface{} {
	return pump.GlucoseUnits()
}

func model(pump *medtronic.Pump, _ Arguments) interface{} {
	return pump.Model()
}

func pumpID(pump *medtronic.Pump, _ Arguments) interface{} {
	return pump.PumpID()
}

func reservoir(pump *medtronic.Pump, _ Arguments) interface{} {
	return pump.Reservoir()
}

func resume(pump *medtronic.Pump, _ Arguments) interface{} {
	log.Printf("resuming pump")
	pump.Suspend(false)
	return nil
}

func rssi(pump *medtronic.Pump, _ Arguments) interface{} {
	return pump.RSSI()
}

func sensitivities(pump *medtronic.Pump, _ Arguments) interface{} {
	return pump.InsulinSensitivities()
}

func setClock(pump *medtronic.Pump, args Arguments) interface{} {
	var t time.Time
	t = parseTime(args["time"].(string))
	log.Printf("setting pump clock to %s", t.Format(medtronic.UserTimeLayout))
	pump.SetClock(t)
	return nil
}

func parseTime(date string) time.Time {
	if date == "now" {
		return time.Now()
	}
	t, err := time.ParseInLocation(medtronic.UserTimeLayout, date, time.Local)
	if err != nil {
		cmdError("setclock", "YYYY-MM-DD HH:MM:SS (or \"now\")", err)
	}
	return t
}

func setMaxBolus(pump *medtronic.Pump, args Arguments) interface{} {
	f, err := args.Float("units")
	if err != nil {
		cmdError("setmaxbolus", "units", err)
	}
	amount := medtronic.Insulin(1000.0*f + 0.5)
	log.Printf("setting max bolus to %v units", amount)
	pump.SetMaxBolus(amount)
	return nil
}

func setTempBasal(pump *medtronic.Pump, args Arguments) interface{} {
	minutes, err := args.Int("duration")
	if err != nil {
		setTempBasalUsage(err)
	}
	duration := time.Duration(minutes) * time.Minute
	f, err := args.Float("rate")
	if err != nil {
		setTempBasalUsage(err)
	}
	temp, err := args.String("temp")
	if err != nil {
		setTempBasalUsage(err)
	}
	switch temp {
	case "absolute":
		rate := medtronic.Insulin(1000.0*f + 0.5)
		log.Printf("setting temporary basal of %v units/hour for %d minutes", rate, minutes)
		pump.SetAbsoluteTempBasal(duration, rate)
	case "percent":
		percent := int(f + 0.5)
		log.Printf("setting temporary basal of %d%% for %d minutes", percent, minutes)
		pump.SetPercentTempBasal(duration, percent)
	default:
		setTempBasalUsage(fmt.Errorf("unknown temp basal type %q", temp))
	}
	return nil
}

func setTempBasalUsage(err error) {
	cmdError("settempbasal", "temp rate duration", err)
}

func settings(pump *medtronic.Pump, _ Arguments) interface{} {
	return pump.Settings()
}

func status(pump *medtronic.Pump, _ Arguments) interface{} {
	return pump.Status()
}

func suspend(pump *medtronic.Pump, _ Arguments) interface{} {
	log.Printf("suspending pump")
	pump.Suspend(true)
	return nil
}

func targets(pump *medtronic.Pump, _ Arguments) interface{} {
	return pump.GlucoseTargets()
}

func tempBasal(pump *medtronic.Pump, _ Arguments) interface{} {
	return pump.TempBasal()
}

func wakeup(pump *medtronic.Pump, _ Arguments) interface{} {
	// pump.Wakeup has already been called
	return nil
}
