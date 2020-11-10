package wrk

import (
	"fmt"

	"github.com/davidkarolyi/server-comparison/benchmark/wrk/types"
	"github.com/gofiber/fiber/v2"
)

var runner = newWRKRunner()

// StartServer will run the wrk benchmarking server
func StartServer() error {
	err := runner.Init()
	if err != nil {
		return err
	}

	app := fiber.New()

	app.Get("/check", check)
	app.Post("/benchmark", benchmark)

	app.Listen(":8080")
	return nil
}

func check(c *fiber.Ctx) error {
	return c.SendStatus(200)
}

func benchmark(c *fiber.Ctx) error {
	params := &types.BenchmarkParams{}
	err := c.BodyParser(params)
	if err != nil {
		return err
	}

	fmt.Printf("Running benchmark: %+v\n", params)
	result, err := runner.RunBenchmark(c.Context(), params)
	if err != nil {
		fmt.Printf("Benchmark failed: '%s'.\n", err)
		return err
	}
	fmt.Println(result.RawOutput)

	return c.JSON(result)
}
