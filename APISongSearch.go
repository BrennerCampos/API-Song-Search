package main

import (
	"encoding/json"
	"fmt"
	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// embedded set of data/structs
type MusicData struct {
	Error    bool `json:"error"`
	Element  widget.Clickable
	Response struct {
		Results []struct { // results is an array struct
			ID       int    `json:"id"`
			Name     string `json:"name"`
			SongName string
		} `json:"results"`
	} `json:"response"`
}

type RelatedData struct {
	Error bool `json:"error"`
	Element  widget.Clickable
	Response struct {
		Similar []struct {
			ID		int `json:"ID"`							// song internal identifier
			Name	string `json:"artist_name"`				// artist name
			SongName	string `json:"song_name"`			// song name
			Lyrics		string `json:"lyrics"`				// song lyrics
			ArtistURL	string `json:"artist_url"`			// artist url
			SongURL		string `json:"song_url"`			// song url
			IndexID		int `json:"index_id"`			// index internal identifier
			PercentSimilar	float32 `json:"percentage"`		//similarity percentage
		} `json:"similarity_list"`
	} `json:"response"`
}

type SongList struct {
	list layout.List
	Items []MusicData
	selected int
}

type RelatedList struct {
	list layout.List
	Items []RelatedData
	selected int
}

var (
	listControl SongList
	sublistControl RelatedList
	appTheme *material.Theme
	songSlice [][]string
	similarSlice [][]string
	topLabel  = "Super Spectacular Song Search"
	editor     = new(widget.Editor)
	lineEditor = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
)

func main() {

	setupList(pullData())
	go startApp()
	app.Main()
}

func setupList(data MusicData) {
	editor.SetText("Write a name of a song or artist: ")
	songSlice = make([][]string, len(data.Response.Results))

	for i := range songSlice {
		songSlice[i] = make([]string, len(data.Response.Results))
	}

	songs := data

	for i := 0; i < len(songs.Response.Results); i++ {
		listControl.Items = append(listControl.Items, songs)
	}
	// making our 2D slice array for storing split results
	for i := 0; i < len(data.Response.Results) ; i++ {

		newString := strings.Split(data.Response.Results[i].Name, "-")

		artistName := strings.Title(strings.ToLower(strings.Trim(newString[0], " ")))
		songName := strings.Title(strings.ToLower(strings.Trim(newString[1], " ")))

		songSlice[i][0] = artistName
		songSlice[i][1] = songName
		songSlice[i][2] = strconv.Itoa(data.Response.Results[i].ID)
		data.Response.Results[i].SongName = songName
	}
	listControl.list.Axis = layout.Vertical
}

func setupSublist(data RelatedData){
	similarSlice = make([][]string, len(data.Response.Similar))

	for i := range songSlice {
		similarSlice[i] = make([]string, len(data.Response.Similar))
	}

	songs := data

	for i := 0; i < len(songs.Response.Similar); i++ {
		sublistControl.Items = append(sublistControl.Items, songs)
	}


	fmt.Println(similarSlice[0])
	//newString := strings.Split(data.Response.Similar[0].SongURL, "-")
	// making our 2D slice array for storing split results
	//for i := 0; i < len(data.Response.Similar) ; i++ {

	//newString := strings.Split(data.Response.Similar[i].Name, "-")

	//artistName := strings.Title(strings.ToLower(strings.Trim(newString[0], " ")))
	//songLyrics := strings.Title(strings.ToLower(strings.Trim(newString[4], " ")))
	//	songName := strings.Title(strings.ToLower(strings.Trim(newString[6], " ")))
	//	similarSlice[i][0] = "test artist_name"//artistName
	///	similarSlice[i][1] = data.Response.Similar[i].ArtistURL
	//	similarSlice[i][2] = strconv.Itoa(data.Response.Similar[i].ID)
	//similarSlice[i][3] = strconv.Itoa(data.Response.Similar[i].IndexID)
	//similarSlice[i][4] = songLyrics
	//similarSlice[i][5] = strconv.FormatFloat(float64(data.Response.Similar[i].PercentSimilar), 'f',-1,32)
	//similarSlice[i][6] = "test"	//songName
	//similarSlice[i][7] = data.Response.Similar[i].SongURL
	//data.Response.Similar[i].SongName = "test2" //songName
	//}
	sublistControl.list.Axis = layout.Vertical
}

func pullData() MusicData {
	var textName string // input from user
	var data MusicData  // making API within main


	fmt.Println("Search for songs! Enter a keyword: 	(\"0\" or \"q\" to quit) ")

	// reads in input
	_, err := fmt.Scan(&textName)
	if err != nil {
		log.Println(err)
	}

	// ways to exit program
	if textName == "0" || textName == "q" || textName == "Q" {
		fmt.Println("Program exiting. Bye!")
		os.Exit(0)
	}


	// https://searchly.asuarez.dev/api/v1/similarity/by_song?song_id=9214   <--- similarity API

	// concatenating user's input to our search API
	response, err := http.Get("https://searchly.asuarez.dev/api/v1/song/search?query=" + textName)
	if err != nil {
		log.Println(err)
	}

	// anonymous function to defer
	defer func() {
		err := response.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// taking in API's data
	rawData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}

	// 'decoding' the data
	err = json.Unmarshal(rawData, &data)
	if err != nil { // & refers to a pointer
		fmt.Println(err)
	}

	//prettyData := data.Response.Results
	CurrentSong := data.Response.Results
	counter := 0

	// printing results
	if len(CurrentSong) <= 0 {
		fmt.Println("Uh-oh, \"" + textName + "\" not found in database. Try again!")
	} else {
		for i := 0; i < len(CurrentSong); i++ {
			data.Response.Results[i].Name = strings.Title(strings.ToLower(data.Response.Results[i].Name))		// making it 'pretty'
			counter++
		}
	}
	fmt.Println()
	fmt.Println("Done! Found " + strconv.Itoa(counter) + " song(s) in total.")

	return data
}

func pullRelated(selectedID int) RelatedData {
	var relatedData RelatedData

	fmt.Println(strconv.Itoa(selectedID))	//debugging incoming selectedID

	//    <--- similarity API
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	subResponse, err := client.Get("https://searchly.asuarez.dev/api/v1/similarity/by_song?song_id="+strconv.Itoa(selectedID))  //+strconv.Itoa(selectedID)
	if err != nil {
		log.Fatal(err)
	}

	/*
		defer func() {
			err := subResponse.Body.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()
		//fmt.Println(response)


	*/
	rawSubdata, err := ioutil.ReadAll(subResponse.Body)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println(string(rawSubdata))

	// 'decoding' the data
	err = json.Unmarshal(rawSubdata, &relatedData)
	if err != nil { // & refers to a pointer
		fmt.Println(err)
	}

	//fmt.Println(&relatedData)

	//prettyData := data.Response.Results

	currentSong := relatedData.Response.Similar
	counter := 0

	//fmt.Println(len(currentSong))
	if len(currentSong) <= 0 {
		fmt.Println("Uh-oh, has no similar songs in the database. Try again!")
	} else {
		for i := 0; i < len(currentSong); i++ {
			relatedData.Response.Similar[i].Name = strings.Title(strings.ToLower(relatedData.Response.Similar[i].Name)) // making it 'pretty'
			counter++
		}
	}

	fmt.Println()
	fmt.Println("Done! Found " + strconv.Itoa(counter) + " similar song(s) in total.")
	//log.Printf(string(rawSubdata))

	return relatedData
	/*
		// concatenating user's input to our search API
		response, err := http.Get("https://searchly.asuarez.dev/api/v1/similarity/by_song?song_id="+strconv.Itoa(selectedID))  // + ID
		if err != nil {
			log.Println(err)
		}

	*/
	/*
		// anonymous function to defer
		defer func() {
			err := response.Body.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()




		// taking in API's data
		rawData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Println(err)
		}


	*/
}

func startApp(){
	defer os.Exit(0) //if we leave this function then  exit with success
	mainWindow := app.NewWindow(app.Size(unit.Dp(600),unit.Dp(600)))		// makes new window


	err := mainEventLoop(mainWindow)	// goes into the main event loop with new window we created
	if err != nil{
		log.Fatal(err)
	}
}


func mainEventLoop(mainWindow *app.Window)(err error){		//takes in a window in the form of a pointer to app.Window

	appTheme = material.NewTheme(gofont.Collection())		//creates our theme

	var operationsQ op.Ops									//creates our operations variable, a way for the program to
	//specify what to display and how to handle events

	for {	//for ever loop
		event := <- mainWindow.Events() //read from the events channel, will wait till there is an event if none
		switch eventType := event.(type){		// figuring out what type of event happened with a switch/case implementation

		case system.DestroyEvent: // so the user closed the window
			return eventType.Err

		case system.ClipboardEvent:
			lineEditor.SetText("event.Text")	//test string, remove quotes

		case system.FrameEvent: //time to draw the window
			graphicsContext := layout.NewContext(&operationsQ, eventType) // creates our graphics context (gtx)
			drawGUI(graphicsContext, appTheme)
			eventType.Frame(graphicsContext.Ops)	// this updates the display of our eventType through our gtx operations
		}		// type FrameEvent.Frame function
	}
}



func drawGUI(gContext layout.Context, theme *material.Theme)layout.Dimensions{	// draws body, or "value" part of the map
	retLayout := layout.Flex{Axis: layout.Horizontal, Alignment: layout.Start}.Layout(gContext, //now we begin building the layout tree toplevel is flex
		layout.Rigid(drawList(gContext, theme)),
		layout.Flexed(0.8, drawSublist(gContext, theme)),
	)
	return retLayout
}


// layout for the top
func drawList(gContext layout.Context, theme *material.Theme)layout.Widget{
	return func(gContext layout.Context) layout.Dimensions { // anonymous func mapping from context to dimensions
		return listControl.list.Layout(gContext, len(listControl.Items),selectItem) //the listItem type is a function from context and int to dimensions
	}
}

func drawSublist(gContext layout.Context, theme *material.Theme)layout.Widget{
	return func(gContext layout.Context) layout.Dimensions { // anonymous func mapping from context to dimensions
		return sublistControl.list.Layout(gContext, len(sublistControl.Items),subselectItem) //the listItem type is a function from context and int to dimensions
	}
}

// a function for what to do when an item gets selected	 --taken as a parameter of list.Layout
func selectItem(graphicsContext layout.Context, selectedItem int) layout.Dimensions{

	for _, e := range lineEditor.Events() {
		if e, ok := e.(widget.SubmitEvent); ok {
			topLabel = e.Text
			lineEditor.SetText("")
		}
	}

	userSelection := &listControl.Items[selectedItem]

	if userSelection.Element.Clicked() {
		listControl.selected = selectedItem
		//fmt.Println("Made it past pull data...")
		//fmt.Println(subData.Response.Similar[selectedItem].SongName)
		setupSublist(pullRelated(userSelection.Response.Results[selectedItem].ID)) //userSelection.Response.Results[selectedItem].ID
	}
	var itemHeight int

	//the layout.Stack.Layout function takes a context followed by possibly many StackChild Structs, each of which must be created
	//by using either layout.Expanded, or layout.Stacked. In either case the parameter is a function for context to dimensions
	return layout.Flex{Alignment: layout.Middle, Axis: layout.Horizontal}.Layout(graphicsContext,	//creates a stack of our elements, in this case a button and selection bar

		// first child							-- stacks our 'buttons' with song artists & names
		layout.Rigid(
			func(gtx layout.Context) layout.Dimensions { 				// 1st anonymous function so we can use userSelection
				dimensions := material.Clickable(gtx, &userSelection.Element,
					func(gtx layout.Context) layout.Dimensions { 		// 2nd anonymous function
						dim := layout.Dimensions{}
						dim = layout.UniformInset(unit.Sp(7)).
							Layout(gtx, material.Body1(appTheme, userSelection.Response.Results[selectedItem].SongName).Layout)
						return dim
					})
				itemHeight = dimensions.Size.Y
				//itemLength = dimensions.Size.X
				return dimensions
			}),
		//end of the first child

		// second  child,			-- paints little rectangle 'selected' bar
		layout.Rigid(
			func(gtx layout.Context) layout.Dimensions { //another one of those 'glorious anonymous functions
				if listControl.selected != selectedItem{
					return layout.Dimensions{} //if not selected - don't do anything special
				}
				paint.ColorOp{Color: color.RGBA{R: 0x80, G: 0x80, A: 0xDE}}.Add(gtx.Ops)// Adds the Primary color from our theme into our gtx operations

				highlightWidth:= gtx.Px(unit.Dp(9)) //lets make it 4 device independent pixels

				paint.PaintOp{Rect: f32.Rectangle{ //paint a rectangle using 32 bit floats
					Max: f32.Point{
						X: float32(highlightWidth),
						Y: float32(itemHeight),
					}}}.Add(gtx.Ops)	// takes all of this in the anonymous function and add it to our gtx operations
				return layout.Dimensions{Size: image.Point{X: highlightWidth, Y: itemHeight}}
			}),
	)
}

func subselectItem(graphicsContext layout.Context, selectedItem int) layout.Dimensions {

	userSubselection := &sublistControl.Items[selectedItem]

	if userSubselection.Element.Clicked() {
		sublistControl.selected = selectedItem
	}
	var itemHeight int

	//the layout.Stack.Layout function takes a context followed by possibly many StackChild Structs, each of which must be created
	//by using either layout.Expanded, or layout.Stacked. In either case the parameter is a function for context to dimensions
	return layout.Stack{Alignment: layout.NW}.Layout(graphicsContext, //creates a stack of our elements, in this case a button and selection bar

		// first child							-- stacks our 'buttons' with song artists & names
		layout.Stacked(
			func(gtx layout.Context) layout.Dimensions { // 1st anonymous function so we can use userSelection
				dimensions := material.Clickable(gtx, &userSubselection.Element,
					func(gtx layout.Context) layout.Dimensions { // 2nd anonymous function
						dim := layout.Dimensions{}
						dim = layout.UniformInset(unit.Sp(7)).
							Layout(gtx, material.Body1(appTheme, userSubselection.Response.Similar[selectedItem].SongName).Layout)
						return dim
					})
				itemHeight = dimensions.Size.Y
				return dimensions
			}),
		//end of the first child

		// second  child,			-- paints little rectangle 'selected' bar
		layout.Stacked(
			func(gtx layout.Context) layout.Dimensions { //another one of those 'glorious anonymous functions
				if sublistControl.selected != selectedItem {
					return layout.Dimensions{} //if not selected - don't do anything special
				}
				paint.ColorOp{Color: color.RGBA{A: 0xDE}}.Add(gtx.Ops) // Adds the Primary color from our theme into our gtx operations
				highlightWidth := gtx.Px(unit.Dp(6)) //lets make it 4 device independent pixels
				paint.PaintOp{Rect: f32.Rectangle{ //paint a rectangle using 32 bit floats
					Max: f32.Point{
						X: float32(highlightWidth),
						Y: float32(itemHeight),
					}}}.Add(gtx.Ops) // takes all of this in the anonymous function and add it to our gtx operations
				return layout.Dimensions{Size: image.Point{X: highlightWidth, Y: itemHeight}}
			}),
	)
}


/*
// Draws The song artist
func drawDisplay(gContext layout.Context, theme *material.Theme)layout.Widget{ //layout.Widget is a function from context to dimensions
	return func (ctx layout.Context) layout.Dimensions {
		displayText := material.H4(theme, "Artist: \n"+songSlice[listControl.selected][0]+"\nID: "+songSlice[listControl.selected][2])  // outputs the artist name for the selected item
		return layout.E.Layout(ctx, displayText.Layout)
	}
}
*/