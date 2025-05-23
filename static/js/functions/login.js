import { clientPage,loadPosts } from "./client.js";
import { connect , fetchUserName} from "./wb.js";

const divLogin = `<img style="margin:9px 0 0 2%;display: none" id="img_close" class="img_close" src="/badr/img/icon_close.png" alt="ssssssssssss">
     <div class="form-container sign-up">
       <div class="form" >
         <h1>Create Account</h1>

         <span>or use your email for registeration</span>
         <input id="FirtsName" type="text" placeholder="Firts Name" />
         <input id="LastName" type="text" placeholder="Last Name" />
         <input id="Nickname" type="text" placeholder="Nickname" />
         <input id="age" type="number" placeholder="Age" />
        <div class="gender-container">
        <input type="radio" id="man" name="Gender" value="man">
        <label for="man">Man</label>

        <input type="radio" id="women" name="Gender" value="women">
        <label for="women">Women</label>
       </div>


         <input id="inputMail" type="email" placeholder="Email" />
         <input id="inputPassword" type="password" placeholder="Password" />
         <button id="regesterform">Sign Up</button>
         
       </div>
     </div>
     <div class="form-container sign-in">
       <div class="form">
         <h1>Sign In</h1>
        
         <span>or use your email password</span>
         <input id="loginInputName" type="email" placeholder="Email or username" />
         <input id="loginInputPassword" type="password" placeholder="Password" />
         <a href="#">Forget Your Password?</a>
         <button id="loginform">Sign In</button>
       </div>
     </div>
     <div class="toggle-container">
       <div class="toggle">
         <div class="toggle-panel toggle-left">
           <h1>Welcome Back!</h1>
           <p>Enter your personal details to use all of site features</p>
           <button id="login">Sign In</button>
         </div>
         <div class="toggle-panel toggle-right">
           <h1>Hello, Friend!</h1>
           <p>
             Register with your personal details to use all of site features
           </p>
           <button id="register">Sign Up</button>
         </div>
       </div>
     </div>`

function loginPage(){
    document.getElementById("container").innerHTML = divLogin
    const container = document.getElementById("container");
    const registerBtn = document.getElementById("register");
    const loginBtn = document.getElementById("login");
    const DivForum = document.getElementById("loginToggle");
    const CloseForum = document.getElementById("img_close");
    const chat = document.getElementById("chatapp");
    document.getElementById("container1").style.display = "none"
    
    logoutButton.style.display = "none"
    chat.style.display = "none"

    DivForum.style.display = "block"
  
    registerBtn.addEventListener("click", () => {
      container.classList.add("active");
    });

    loginBtn.addEventListener("click", () => {
      container.classList.remove("active");
    });

    DivForum.addEventListener("click", () => {
      container.style.display = "inline-block"
        document.getElementById("img_close").style.display = "inline-block"

    });

    CloseForum.addEventListener("click", () => {
      container.style.display = "none"
        document.getElementById("img_close").style.display = "none"

    });
}
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

function loginHundler(){
    document.getElementById("regesterform").addEventListener("click",async function (){
        const FirtsName = document.getElementById("FirtsName").value;
        const LastName = document.getElementById("LastName").value;
        const Nickname = document.getElementById("Nickname").value;
        const age = document.getElementById("age").value;
        const inputMail = document.getElementById("inputMail").value;
        const inputPassword = document.getElementById("inputPassword").value;
        let gender1;
        const gender = document.querySelectorAll('input[name="Gender"]');
        gender.forEach((radio) => {
          if (radio.checked) {
            gender1 = radio.value;
          }
      });

        
        const formData = new URLSearchParams();
        formData.append("FirtsName",FirtsName)
        formData.append("LastName",LastName)
        formData.append("username",Nickname)
        formData.append("age",age)
        formData.append("gender",gender1)
        formData.append("email",inputMail)
        formData.append("password",inputPassword)
      
      
        try{
          const response = await fetch("/regester",{
            method: "POST",
            body:formData
          })
          const Json = await response.json();
          console.log(Json);
          
          if (!response.ok){
            alert(Json.message)
            
          }else{
            alert(Json.message)
            document.getElementById("container").classList.remove("active")
          }
        }catch (error){
          
        }
        
      })
      
      document.getElementById("loginform").addEventListener("click",async function (){
        const inputName = document.getElementById("loginInputName").value;
        const inputPassword = document.getElementById("loginInputPassword").value;
      
        const formData = new URLSearchParams();
        formData.append("email",inputName)
        formData.append("password",inputPassword)
      
        let response = await fetch("/login",{
            method: "POST",
            body:formData
        })
        
        
        const Json = await response.json();

        if (!response.ok){
          alert(Json.message)
        }else{
          alert(Json.message)
          clientPage()
          loadPosts()
          connect(Json.username)
          fetchUserName(Json.username)
        }   
      })
}

export {loginPage,loginHundler}