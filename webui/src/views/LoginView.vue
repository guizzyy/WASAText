<script>
  import ErrorMsg from "../components/ErrorMsg.vue";

  export default {
    components: {ErrorMsg},
    data: function () {
      return {
        error: null,
        photo: "",
        username: "",
        ID: 0,
        message: ""
      }
    },
    methods:{
      async doLogin(e){
        e.preventDefault();
        if (this.username === ""){
          this.error = "Username is required";
        } else {
          this.error = null;
          try {
            let response = await this.$axios.post("/session", {username: this.username})
            sessionStorage.setItem("ID", response.data.user.id);
            sessionStorage.setItem("username", response.data.user.username);
            sessionStorage.setItem("photo", response.data.user.photo);
            sessionStorage.setItem("message", response.data.message);
            this.$router.push({ path : "/conversations" })
          } catch (e) {
            if (e.response && e.response.status === 400) {
              this.error = "Invalid username (it must be between 3 and 16 characters.)";
            } else if (e.response && e.response.status === 500) {
              this.error = "Server Error, please try again later.";
            } else {
              this.error = e.toString();
            }
          }
        }
        setTimeout(() => {
          this.error = null;
        }, 2500)
      }
    }
  }
</script>

<template>

  <div class="d-flex position-relative">
    <div class="d-flex position-absolute top-0 end-0 mt-3" style="padding-right: 10px">
      <ErrorMsg v-if="error" :msg="error"></ErrorMsg>
    </div>
  </div>

  <div class="d-flex justify-content-center position-absolute" style="top: 27%; width: 100%; height: 100%;">
    <div class="justify-content-between flex-wrap flex-md-nowrap align-items-center">
      <h2 class="h2 text-center">Welcome to WASAText</h2>
    </div>
  </div>

  <div class="d-flex justify-content-center position-absolute" style="top: 40%; width: 100%; height: 100%;">
    <h2 class="h2 text-center" v-if="username"> {{ username }}</h2>
  </div>

  <div class="d-flex justify-content-center position-absolute" style="top: 52%; left: 0; width: 100%; height: 100%; padding-top: 1.5rem">
    <form @submit="doLogin" class="mt-6">
      <div class="flex items-center justify-center min-h-screen">
        <input id="username-given" v-model="username" type="text" placeholder="Enter your username" autocomplete="off" maxlength="16"
               class="w-full p-3 rounded-md text-black border border-gray-300 focus:outline-none focus:ring-2 focus:ring-blue-500 form-control">
        <div class="text-center" style="padding-top: 1.5rem">
          <button type="submit">
              <span>Start Chatting</span>
          </button>
        </div>
      </div>
    </form>
  </div>


</template>

<style>

body, html {
  height: 100%;
  margin: 0;
  padding: 0;
}

.h2 {
  color: snow;
  display: flex;
  justify-content: center;
  position: absolute;
  width: 100%;
  height: 100%;
  left: 0;
  font-size: 5rem;
  font-weight: bold;
  text-shadow: 2px 2px 3px #5c636a;
}

.site-name h1 {
  font-size: 6rem;
}

form {
  width: 100%;
  max-width: 400px;
}

</style>