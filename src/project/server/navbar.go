package server

import (
    "html/template"

    log "project/logging"
)

type Menu struct {
    Navbar bool
    Title string
    Site template.URL
    Text string
    LoggedIn bool
    LeftItems []MenuItem
    RightItems []MenuItem
}

type MenuItem struct {
    Active bool
    Link template.URL
    Text string
    Dropdown bool
    DropItems []DropItem
    Id string
}

type DropItem struct {
    Link template.URL
    Text string
}

func ImageNav(user_id string, topsite string) Menu {
    log.Printf("ImageNav: %s", topsite)
    loggedIn := user_id != ""

    submenu := []DropItem{{"#", "Page1"}, {"#", "Page2"}}
    navbar := Menu{true, "Sample Title", template.URL(topsite), "MySite", loggedIn,
        []MenuItem{
            {Active:true, Link: template.URL(topsite), Text:"Home"},
            {Link:"#", Text:"About"},
            {Dropdown:true, Link:"#", Text:"LeftMenu", DropItems:submenu, Id:"my-left-dropmenu"},
        }, []MenuItem{
            {Dropdown:true, Link:"#", Text:"RightMenu", DropItems:submenu, Id:"my-right-dropmenu"},
        },
    }

    if loggedIn {
        navbar.RightItems = append(navbar.RightItems, MenuItem{Link:"#", Text:"Account"})
    }
    return navbar
}
