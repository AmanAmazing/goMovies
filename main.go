package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
    "github.com/rivo/tview"
    "github.com/gdamore/tcell/v2"
)


func init(){
    err := godotenv.Load(".env")
    if err !=nil{
        log.Fatal("Error loading .env file. Please ensure the .env file is setup correctly")
    }
}

var name string 
var help bool
var api_key string
var list *tview.List
var textDisplay *tview.TextView
var flex = tview.NewFlex()
var pages = tview.NewPages()

func dealWithFlags(){
    flag.StringVar(&name,"name","","specify name for movie/show, REQUIRED")
    flag.BoolVar(&help,"help",false,"Display usage information")
    flag.Parse()
    if help {
        flag.PrintDefaults()
        os.Exit(0)
    }
    if name == ""{
        flag.PrintDefaults()
        log.Fatal("input path is required")
    }
}

func menu(app *tview.Application,moviesList []Search) *tview.List{
    list = tview.NewList()
    for i, p := range moviesList{
        list.AddItem(p.Title,p.ImdbID,rune(49+i),nil)
    }

    list.AddItem("Quit","Press to exit",'q',func(){
        app.Stop()
    })

    return list 
}

func setResponse(index int, mainText string, secondaryText string, shortcut rune){
    movie := findById(secondaryText) 
    //list.SetItemText(index,mainText,movie.Plot)
    setTextDisplay(&movie)    
}

func setTextDisplay(m *Search){
    textDisplay.Clear()
    text := m.Runtime + "\n" + m.Director + "\n" + m.Plot
    textDisplay.SetText(text).
        SetTextAlign(tview.AlignCenter).
        SetTextColor(tcell.ColorGreen).SetBorderPadding(5,5,5,5)
}


func searchByName() ResponseSearch{
    api_key = os.Getenv("API_URL")
    searchUrl := fmt.Sprint(api_key+"&s="+name)
    
    res, err := http.Get(searchUrl)
    if err != nil {
        log.Fatal(err)
    }
    defer res.Body.Close()
    body, err := io.ReadAll(res.Body)
    if err != nil {
        log.Fatal(err)
    }

    var searchResults ResponseSearch

    if err := json.Unmarshal(body, &searchResults); err != nil {
        log.Printf("Could not unmarshal responseBytes: %v",err)
    }

    if searchResults.Response == "False"{
        fmt.Println(searchResults.Error)
    }
    return searchResults
}



func findById(id string) Search{
    searchUrl := fmt.Sprint(api_key+"&i="+ id) 

    res, err := http.Get(searchUrl)
    if err != nil {
        log.Fatal(err)
    }
    defer res.Body.Close()
    body, err := io.ReadAll(res.Body)
    if err != nil {
        log.Fatal(err)
    }

    var movieResults Search

    if err := json.Unmarshal(body, &movieResults); err != nil {
        log.Printf("Could not unmarshal responseBytes: %v",err)
    }
    
    return movieResults
    
}

func main() {
    dealWithFlags()

    results := searchByName()

    app := tview.NewApplication()
    list := menu(app, results.Search)
    //list.SetChangedFunc(setResponse)
    list.SetSelectedFunc(setResponse) 
    // text box 
    textDisplay = tview.NewTextView()

    // layout
    flex.SetDirection(tview.FlexRow).
    AddItem(tview.NewFlex().
    AddItem(list,0,1,true).
    AddItem(textDisplay,0,4,false),0,6,false)
    pages.AddPage("filmsList",flex,true,true)
    if err := app.SetRoot(pages,true).EnableMouse(true).Run(); err !=nil{
        panic(err)
    }
}

type ResponseSearch struct{
    Search []Search `json:"Search"` 
    Total string `json:"totalResults"` 
    Response string `json:"Response"`
    Error string `json:"Error"`
} 

type Search struct {
    Title string `json:"Title"`
    Year string `json:"Year"`
    Runtime string `json:"Runtime"`
    Genre string `json:"Genre"`
    Director string `json:"Director"`
    Plot string `json:"Plot"`
    ImdbID string `json:"imdbID"`
    Type string `json:"Type"`
    Poster string `json:"Poster"`
    Ratings []Rating
}

type Rating struct {
    Source string `json:"Source"`
    Value string `json:"Value"`
}
