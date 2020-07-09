package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	wakabox "github.com/YouEclipse/waka-box-go/pkg"
	"github.com/google/go-github/github"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load(".env")
}

func main() {
	err := godotenv.Load(".env")

	wakaAPIKey := os.Getenv("WAKATIME_API_KEY")
	ghToken := os.Getenv("GH_TOKEN")
	ghUsername := os.Getenv("GH_USER")
	gistID := os.Getenv("GIST_ID")

	style := wakabox.BoxStyle{
		BarStyle:  os.Getenv("GIST_BARSTYLE"),
		BarLength: os.Getenv("GIST_BARLENGTH"),
		TimeStyle: os.Getenv("GIST_TIMESTYLE"),
	}

	dryRun := os.Getenv("DRY_RUN")

	box := wakabox.NewBox(wakaAPIKey, ghUsername, ghToken, style)

	lines, err := box.GetStats(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println(strings.Join(lines, "\n"))
	if dryRun == "TRUE" {
		return
	}
	ctx := context.Background()

	filename := "📊 Weekly development breakdown"
	gist, err := box.GetGist(ctx, gistID)
	if err != nil {
		panic(err)
	}

	f := gist.Files[github.GistFilename(filename)]

	f.Content = github.String(strings.Join(lines, "\n"))
	gist.Files[github.GistFilename(filename)] = f
	err = box.UpdateGist(ctx, gistID, gist)
	if err != nil {
		panic(err)
	}
}
