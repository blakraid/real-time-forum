import {loginPage, loginHundler } from './functions/login.js'
import { clientPage, loadPosts } from './functions/client.js'
import {connect, fetchUserName} from './functions/wb.js'

document.addEventListener("DOMContentLoaded", async () =>{
    const response = await fetch("/check-session")
    const reJson = await response.json()
    
    if (!reJson.loggedIn){
        loginPage()
        loginHundler()
    }else {
        document.getElementById("homepage").style.display = "block"
        await fetchUserName(reJson.username)
        connect(reJson.username)
        clientPage()
        loadPosts()
    }
})
