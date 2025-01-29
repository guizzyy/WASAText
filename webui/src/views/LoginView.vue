<script>
  import ErrorMsg from "../components/ErrorMsg.vue";

  export default {
    components: {ErrorMsg},
    data: function () {
      return {
        error: null,
        photo: "",
        username: "",
        ID: 0
      }
    },
    methods:{
      async doLogin(){
        if (this.username === ""){
          this.error = "Username is required";
        } else {
          this.error = null;
          try {
            let response = await this.$axios.post("/session", {username: this.username})
            sessionStorage.setItem("ID", response.data.id);
            sessionStorage.setItem("username", response.data.username);
            sessionStorage.setItem("photo", response.data.photo);
            this.$router.push({ path : "/conversations" })
          } catch (e) {
            if (e.response && e.response.status === 400) {
              this.error = "Invalid username (it must be between 3 and 16 characters.)";
            } else if (e.response && e.response.status === 500) {
              this.error = "Server Error, please try again later.";
            } else {
              this.error = e.toString();
            }
            setTimeout(() => {
              this.error = null;
            }, 5000)
          }
        }
      }
    }
  }
</script>

<template>
  <div>
    <ErrorMsg v-if="error" :msg="error"></ErrorMsg>/
  </div>

  <div>
    <div class="site-name" style="top: 15%">
      <h1>WASAText</h1>
    </div>

    <div style="top: 40%; width: 100%; text-align: center; height: 100%">
      <div class="Message title">
        <h2>Welcome to WASAText, {{username}}!</h2>
      </div>
    </div>
  </div>

  <div style="top: 50%; left: 0; width: 100%; height: 100%; padding-top: 1.25rem">
    <form @submit.prevent="doLogin">
      <div>
        <input id="username-given" v-model="username" type="text" placeholder="Enter your username" autocomplete="off"></input>
        <div>
          <button type="submit">
            <!-- INSERT AN IMAGE FOR THE BUTTON -->
            <span>Login</span>
          </button>
        </div>
      </div>
    </form>
  </div>



</template>

<style scoped>

body, html {
  height: 100%;
  margin: 0;
  padding: 0;
}

.site-name {
  display: flex;
  justify-content: center;
  position: absolute;
  width: 100%;
  height: 100%;
  left: 0;
  font-size: 3rem;
  font-weight: bold;
  text-shadow: 2px 2px 5px #b02a37;
}

.site-name h1 {
  font-size: 6rem;
}

form {
  width: 100%;
  max-width: 400px;
}

</style>