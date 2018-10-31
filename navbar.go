package main

import (
    "html/template"
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
    loggedIn := user_id != ""

    submenu := []DropItem{{Link:"#", Text:"Page1"}, {Link:"#", Text:"Page2"}}
    navbar := Menu{Navbar:true, Title:"Sample Title", Site:template.URL(topsite), Text:"MySite", LoggedIn:loggedIn,
        LeftItems:[]MenuItem{
            {Active:true, Link:template.URL(topsite), Text:"Home"},
            {Link:"#", Text:"About"},
            {Dropdown:true, Link:"#", Text:"LeftMenu", DropItems:submenu, Id:"my-left-dropmenu"},
        }, RightItems:[]MenuItem{
            {Dropdown:true, Link:"#", Text:"RightMenu", DropItems:submenu, Id:"my-right-dropmenu"},
        },
    }

    if loggedIn {
        navbar.RightItems = append(navbar.RightItems, MenuItem{Link:template.URL(topsite+"/account"), Text:"Account"})
    }
    return navbar
}*/
