package main

import (
    "fmt"
    "github.com/rivo/tview"
    "github.com/gdamore/tcell/v2"
    "strconv"
)
//Pre:
//Post: it prints a tree where every channel of the exported XML
func printChannels(programs []*Program) *tview.Flex {
	root := tview.NewTreeNode("Channels").
		SetColor(tcell.ColorRed)
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	add := func(target *tview.TreeNode, entry int) {
		if (entry == -1) {
            for i := 0; i < len(programs); i++ {
                node := tview.NewTreeNode(trimFirstThirdRuns(programs[i].Name)).
    				SetReference(i).
    				SetSelectable(true)
    				node.SetColor(tcell.ColorGreen)
    			target.AddChild(node)
            }
        } else {
            node_videos := tview.NewTreeNode("Videos")
            target.AddChild(node_videos)
            for i := 0; i < len(programs[entry].Videos); i++ {
                subnode_video := tview.NewTreeNode("Video "+ strconv.Itoa(i+1))
                video_child := tview.NewTreeNode("Format: " + programs[entry].Videos[i].Format)
                node_videos.AddChild(subnode_video)
                subnode_video.AddChild(video_child)
            }
            node_audios := tview.NewTreeNode("Audios")
            target.AddChild(node_audios)
            for i := 0; i < len(programs[entry].Audios); i++ {
                subnode_audio := tview.NewTreeNode("Audio "+ strconv.Itoa(i+1))
                audio_child1 := tview.NewTreeNode("Language: " + programs[entry].Audios[i].Language)
                audio_child2 := tview.NewTreeNode("Format: " + programs[entry].Audios[i].Format)
                node_audios.AddChild(subnode_audio)
                subnode_audio.AddChild(audio_child1)
                subnode_audio.AddChild(audio_child2)
            }
        }
	}

	// Add all the channels.
	add(root, -1)

	// If a channel was selected, open it.
	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			return // Selecting the root node does nothing.
		}
		children := node.GetChildren()
		if len(children) == 0 {
			add(node, reference.(int))
		} else {
			// Collapse if visible, expand if collapsed.
			node.SetExpanded(!node.IsExpanded())
		}
	})

    main := tview.NewFlex().
        AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
            AddItem(tree,0,10,true), 0, 2, true)
	return main
}

//Pre:
//Post: it prints a tree where every entry is a channel of the desired list
//envery entry has as childs its matches
func printmatches(matches [][]*Program, desiredChList []ChListElement) *tview.Flex {
	root := tview.NewTreeNode("Channels").
		SetColor(tcell.ColorRed)
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	add := func(target *tview.TreeNode, entry int) {
		if (entry == -1) {
            for i := 0; i < len(desiredChList); i++ {
                node := tview.NewTreeNode("[" + fmt.Sprintf("%03d", i+1) + "][" + strconv.Itoa(len(matches[i])) + " match] " +
                desiredChList[i].DisplayName).
    				SetReference(i).
    				SetSelectable(true)

    			if (len(matches[i]) == 1) {
                    node.SetColor(tcell.ColorGreen)
                } else if (len(matches[i]) == 0) {
                    node.SetColor(tcell.ColorRed)
                } else {
                    node.SetColor(tcell.ColorOrange)
                }
    			target.AddChild(node)
            }
        } else {
            for c := 0; c < len(matches[entry]); c++ {
                connection := tview.NewTreeNode(trimFirstThirdRuns(matches[entry][c].Name))
                target.AddChild(connection)
                node_videos := tview.NewTreeNode("Videos")
                connection.AddChild(node_videos)
                for i := 0; i < len(matches[entry][c].Videos); i++ {
                    subnode_video := tview.NewTreeNode("Video "+ strconv.Itoa(i+1))
                    video_child := tview.NewTreeNode("Format: " + matches[entry][c].Videos[i].Format)
                    node_videos.AddChild(subnode_video)
                    subnode_video.AddChild(video_child)
                }
                node_audios := tview.NewTreeNode("Audios")
                connection.AddChild(node_audios)
                for i := 0; i < len(matches[entry][c].Audios); i++ {
                    subnode_audio := tview.NewTreeNode("Audio "+ strconv.Itoa(i+1))
                    audio_child1 := tview.NewTreeNode("Language: " + matches[entry][c].Audios[i].Language)
                    audio_child2 := tview.NewTreeNode("Format: " + matches[entry][c].Audios[i].Format)
                    node_audios.AddChild(subnode_audio)
                    subnode_audio.AddChild(audio_child1)
                    subnode_audio.AddChild(audio_child2)
                }
            }
        }
	}

	// Add all the channels.
	add(root, -1)

	// If a channel was selected, open it.
	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			return // Selecting the root node does nothing.
		}
		children := node.GetChildren()
		if len(children) == 0 {
			add(node, reference.(int))
		} else {
			// Collapse if visible, expand if collapsed.
			node.SetExpanded(!node.IsExpanded())
		}
	})

    main := tview.NewFlex().
        AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
          AddItem(tree,0,10,true), 0, 2, true)
	return main
}

//Pre:
//Post: it prints a tree where every entry is a channel of the desired list
//envery entry has as childs its matches. It only displays the elements without
//matches or the elements with more than a match. It lets user to pick a match from
//all the matches for the channel.
func selectmatches(db *Database, matches [][]*Program, desiredChList []ChListElement, app *tview.Application) *tview.Flex {
	root := tview.NewTreeNode("Channels that have more than a match").
		SetColor(tcell.ColorRed).
        SetSelectable(false)
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

    var matchchannels int = 0
    var emptychannels int = 0

	add := func(target *tview.TreeNode) {
        for i := 0; i < len(desiredChList); i++ {
            if (len(matches[i]) == 0) {
                //desired channel without any match
                if(desiredChList[i].SearchName != "%empty%") {
                    node := tview.NewTreeNode("[" + fmt.Sprintf("%03d", i+1) + "]" +
                    desiredChList[i].DisplayName + " has no matches. Try using SearchName in desired.txt. Stopping...").
                    SetColor(tcell.ColorRed).
                    SetSelectable(false)
                    target.AddChild(node)
                } else {
                    //if the channel has to be empty, generate an empty channel
                    emptychannels++
                }

            } else if (len(matches[i]) > 1) {
                //more than a match. the user needs to pick only a channel
                node := tview.NewTreeNode("[" + fmt.Sprintf("%03d", i+1) + "][" + strconv.Itoa(len(matches[i])) + " match] " +
                desiredChList[i].DisplayName).
    				SetReference(i).
                    SetSelectable(false)
                node.SetColor(tcell.ColorOrange)
    			target.AddChild(node)

                for c := 0; c < len(matches[i]); c++ {
                    //add an int reference. it will be used for identifying it when selecting a match
                    connection := tview.NewTreeNode(trimFirstThirdRuns(matches[i][c].Name)).
                        SetReference((100*i)+(c)).
                        SetSelectable(true).
                        SetColor(tcell.ColorBlue)

                    node.AddChild(connection)

                    node_videos := tview.NewTreeNode("Videos").SetSelectable(false)
                    connection.AddChild(node_videos)
                    for v := 0; v < len(matches[i][c].Videos); v++ {
                        subnode_video := tview.NewTreeNode("Video "+ strconv.Itoa(v+1)).SetSelectable(false)
                        video_child := tview.NewTreeNode("Format: " + matches[i][c].Videos[v].Format).SetSelectable(false)
                        node_videos.AddChild(subnode_video)
                        subnode_video.AddChild(video_child)
                    }

                    node_audios := tview.NewTreeNode("Audios").SetSelectable(false)
                    connection.AddChild(node_audios)
                    for a := 0; a < len(matches[i][c].Audios); a++ {
                        subnode_audio := tview.NewTreeNode("Audio "+ strconv.Itoa(a+1)).SetSelectable(false)
                        audio_child1 := tview.NewTreeNode("Language: " + matches[i][c].Audios[a].Language).SetSelectable(false)
                        audio_child2 := tview.NewTreeNode("Format: " + matches[i][c].Audios[a].Format).SetSelectable(false)
                        node_audios.AddChild(subnode_audio)
                        subnode_audio.AddChild(audio_child1)
                        subnode_audio.AddChild(audio_child2)
                    }
                }

            } else {
                matchchannels++
                continue }
        }
	}

	// Add all the channels that have matching problems
	add(root)

    var main *tview.Flex
    if (matchchannels+emptychannels != len(desiredChList)) {
        //there are channels to modify. show the tree
        main = tview.NewFlex().
          AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
          AddItem(tree,0,10,true), 0, 2, true)

    	tree.SetSelectedFunc(func(node *tview.TreeNode) {
            //this function is called when selecting an intem in the tree list
            //reference must be the result of i*100+c (matches[i][c])
    		reference := node.GetReference()
    		if reference == nil {
    			panic("no reference!")
    		}

            i := reference.(int)/100
            c := reference.(int)%100

            var selectedCh []*Program
            selectedCh = append(selectedCh,matches[i][c])
            matches[i] = selectedCh

            root.ClearChildren()
            add(root)
    	})
    } else {
        //matching process is done
        //numerate all the channels
        for c := 0; c < len(matches); c++ {
            if (desiredChList[c].DisplayName == "%empty%") {
                numerateChannel(matches[c][0],c+1,desiredChList[c].DisplayName)
            } else {
                res := numerateChannel(matches[c][0],c+1,desiredChList[c].DisplayName)
                if (res<0) {
                    message := "PANIC. Channel " + trimFirstThirdRuns(matches[c][0].Name) +
                      " is already assigned in channel number " + strconv.Itoa(-1*res) +
                      ". Desired channel name was " + desiredChList[c].DisplayName +
                      ", in position number " + strconv.Itoa(c)
                    displayError(app,message)
                }
            }
        }
        modal := tview.NewModal().
		SetText("All the channels have a match! What do you want to do with " +
             "the remaining channels of your channel list?").
		AddButtons([]string{"Delete them", "Keep only radio channels", "Keep them all"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Delete them" {
                cleanup(db, false)
			} else if buttonLabel == "Keep only radio channels" {
                cleanup(db, true)
            }
            exportXML(db,app)
		})
        main = tview.NewFlex().
            AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
                AddItem(modal,0,10,true), 0, 2, true)
    }
    return main
}

func displayError(app *tview.Application, message string) {
    var main *tview.Flex
    modal := tview.NewModal().
    SetText(message).
    AddButtons([]string{"Quit"}).
    SetBackgroundColor(tcell.ColorRed).
    SetDoneFunc(func(buttonIndex int, buttonLabel string) {
        if buttonLabel == "Quit" {
            app.Stop()
        }
    })

    main = tview.NewFlex().
      AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
      AddItem(modal,0,10,true), 0, 2, true)

    update_ui(app,main,":-(")
}

//Pre
//Post: exports the XML and shows a message in the screen
func exportXML(db *Database, app *tview.Application) {
    var main *tview.Flex
    filename:=writeFile(db)
    modal := tview.NewModal().
    SetText("Export success :-)\nFile created: " + filename).
    AddButtons([]string{"Menu", "Quit"}).
    SetBackgroundColor(tcell.ColorGreen).
    SetDoneFunc(func(buttonIndex int, buttonLabel string) {
        if buttonLabel == "Menu" {
            menu(app)
        } else {
            app.Stop()
        }
    })

    main = tview.NewFlex().
      AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
      AddItem(modal,0,10,true), 0, 2, true)

    update_ui(app,main,":-D")
}

//Pre:
//Post: the UI is shown. KeyTab is used for passing to the next screen
func makeChList(db *Database, app *tview.Application, flex *tview.Flex, desiredChList []ChListElement, matches [][]*Program) {
    main := selectmatches(db, matches, desiredChList, app)
    app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Anything handled here will be executed on the main thread
		switch event.Key() {
    		case tcell.KeyEsc:
    			// Exit the application
    			menu(app)
    			return nil
            case tcell.KeyTab:
    			// Exit the application
    			makeChList(db, app, flex, desiredChList, matches)
    			return nil
		}

		return event
	})
    update_ui(app,main,"Pick the channel you want for every entry using Enter key. When you have picked all the channels, use the TAB key to continue")
}


//Pre: flex view initializated, footer_text is a string
//Post: it displays the main flex view with a fixed heather and the footer with its text
func update_ui(app *tview.Application, main *tview.Flex, footer_text string) {
    header := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetTextAlign(tview.AlignCenter).
		SetTextColor(tcell.ColorBlack)
	header.Box.SetBackgroundColor(tcell.ColorWhite)
	fmt.Fprintf(header,"GTC Channel List Maker by osanchezdm - v0.01")

	footer := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetTextColor(tcell.ColorWhite)
	footer.Box.SetBackgroundColor(tcell.ColorGrey)
	fmt.Fprintf(footer,footer_text)

    flex := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(header,1,1,false).
			AddItem(main,0,10,true).
			AddItem(footer, 1, 1, false), 0, 2, false)

    if err := app.SetRoot(flex, true).SetFocus(main).Run(); err != nil {
        panic(err)
    }
}

//Post: it displays the main menu
func menu(app *tview.Application) {
    var main *tview.Flex
    err1,dbptr := getDatabase()
    db := *dbptr

    err2,desiredChList := getDesiredChList()

    if (err1<0) {
        displayError(app,"Error loading Database. Make sure input.xml it's inside the app folder")
        app.Stop()
    } else if (err2<0) {
        displayError(app,"Error loading desired channel list. Make sure desired.txt it's inside the app folder")
        app.Stop()
    } else {
        dbClean(&db,0)
        var programpointers []*Program = program_list(dbptr)

        list := tview.NewList().
        AddItem("List XML channels", "Prints a list of the channels reading the XML file exported by your GTC.", 'l', func() {
                main := printChannels(programpointers)
                update_ui(app,main,"Use the up/down arrows to move through the list. Use the Esc key to return to the main menu.")
            }).
        AddItem("Check desired channels", "Checks if candidates have been found for the desired list.", 'b', func() {
                matches := makeMatchMatrix(programpointers,desiredChList, dbptr)
                main := printmatches(matches, desiredChList)
                update_ui(app,main,"Use the up/down arrows to move through the list. Use the Esc key to return to the main menu.")
            }).
        AddItem("Make a list", "Applicates the channel order and asks whether the user wants to keep or delete the remaining channels.", 'm', func() {
                matches := makeMatchMatrix(programpointers,desiredChList, dbptr)
                makeChList(dbptr,app,main,desiredChList,matches)
            }).
        AddItem("Help", "Please check https://github.com/oscarsanchezdm/channellistmaker", 'h', func() {

            }).
        AddItem("Quit", "Press to exit.\n", 'q', func() {
            app.Stop()
        }).SetSelectedBackgroundColor(tcell.ColorDarkGrey)

        //Main window

        intro := tview.NewTextView().
    		SetDynamicColors(true).
    		SetRegions(true)
    	fmt.Fprintf(intro,"Welcome to GTC Channel List Maker! Please pick an option to continue.")

        main = tview.NewFlex().
            AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
                AddItem(intro,2,2,false).
                AddItem(list,0,10,true), 0, 2, true)

        footer_text := "Press q to quit"
        update_ui(app, main, footer_text)
    }
}

func main() {
    app := tview.NewApplication()
    app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Anything handled here will be executed on the main thread
		switch event.Key() {
    		case tcell.KeyEsc:
    			// Exit the application
    			menu(app)
    			return nil
		}
		return event
	})
    //do here the db load
    //do here the desiredchlist load
    //look here for the matches
    menu(app)
    app.Stop()
}
