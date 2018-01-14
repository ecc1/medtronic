package main

import (
	"log"
	"strconv"
	"time"

	"github.com/ecc1/medtronic"
)

type (
	// Arguments represents the formal and actual parameters for a command.
	Arguments map[string]interface{}

	// Prog represents a function to be executed with the given arguments on the pump.
	Prog func(*medtronic.Pump, Arguments) interface{}

	// Command specifies the function and formal parameters for a command.
	Command struct {
		Cmd    Prog
		Params []string
	}
)

var (
	// TODO: add per-command help
	command = map[string]Command{
		"basal":         cmd(basal),
		"battery":       cmd(battery),
		"bolus":         cmd(bolus, "units"),
		"button":        cmd(button, "keys..."),
		"carbratios":    cmd(carbRatios),
		"carbunits":     cmd(carbUnits),
		"clock":         cmd(clock),
		"execute":       cmd(execute, "command", "arguments..."),
		"glucoseunits":  cmd(glucoseUnits),
		"history":       cmd(history),
		"model":         cmd(model),
		"pumpid":        cmd(pumpID),
		"reservoir":     cmd(reservoir),
		"resume":        cmd(resume),
		"rssi":          cmd(rssi),
		"sensitivities": cmd(sensitivities),
		"setclock":      cmd(setClock, "date", "time"),
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
	return Command{Cmd: prog, Params: params}
}

// TODO: with argument to schedule progs, get schedule at that time

func basal(pump *medtronic.Pump, _ Arguments) interface{} {
	return pump.BasalRates()
}

func battery(pump *medtronic.Pump, _ Arguments) interface{} {
	return pump.Battery()
}

func bolus(pump *medtronic.Pump, args Arguments) interface{} {
	f, err := strconv.ParseFloat(args["units"].(string), 64)
	if err != nil {
		bolusUsage()
	}
	amount := medtronic.Insulin(1000.0*f + 0.5)
	log.Printf("performing bolus of %v units", amount)
	pump.Bolus(amount)
	return nil
}

func bolusUsage() {
	log.Fatal("usage: bolus units")
}

func button(pump *medtronic.Pump, args Arguments) interface{} {
	for _, s := range args["keys..."].([]string) {
		b := parseButton(s)
		log.Printf("pressing %v", b)
		pump.Button(b)
		if pump.Error() != nil {
			log.Fatal(pump.Error())
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
		log.Printf("unknown pump button (%s)", s)
		buttonUsage()
	}
	return b
}

func buttonUsage() {
	log.Fatal("usage: button (b|esc|act|up|down) ...\n")
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
	cmd, err := strconv.ParseUint(args["command"].(string), 16, 8)
	if err != nil {
		log.Printf("%v", err)
		executeUsage()
	}
	var params []byte
	for _, s := range args["arguments..."].([]string) {
		b, err := strconv.ParseUint(s, 16, 8)
		if err != nil {
			log.Printf("%v", err)
			executeUsage()
		}
		params = append(params, byte(b))
	}
	return pump.Execute(medtronic.Command(cmd), params...)
}

func executeUsage() {
	log.Fatal("usage: execute cmd [param ...]")
}

func glucoseUnits(pump *medtronic.Pump, _ Arguments) interface{} {
	return pump.GlucoseUnits()
}

func history(pump *medtronic.Pump, _ Arguments) interface{} {
	return pump.LastHistoryPage()
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
	ds := args["date"].(string)
	ts := args["time"].(string)
	if ds == "now" {
		t = time.Now()
	} else {
		t = parseTime(ds + " " + ts)
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
	log.Fatal("usage: setclock YYYY-MM-DD HH:MM:SS (or \"now\")")
}

func setTempBasal(pump *medtronic.Pump, args Arguments) interface{} {
	minutes, err := strconv.ParseUint(args["duration"].(string), 10, 8)
	if err != nil {
		setTempBasalUsage()
	}
	duration := time.Duration(minutes) * time.Minute
	f, err := strconv.ParseFloat(args["rate"].(string), 64)
	if err != nil {
		setTempBasalUsage()
	}
	temp := args["temp"].(string)
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
		setTempBasalUsage()
	}
	return nil
}

func setTempBasalUsage() {
	log.Fatal("usage: settempbasal temp rate duration")
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
