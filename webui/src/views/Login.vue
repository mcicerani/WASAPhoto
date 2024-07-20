<script>
import api from "@/services/axios"

export default {
  data() {
    return {
      username: '',
    };
  },
  methods: {
    async dologin() {
      try {
        console.log('Logging in with username:', this.username);  // Debugging step

        const response = await api.post('/session', {
          username: this.username,
        });

        const token = response.data.token;
        console.log('Token received:', token);  // Debugging step

        const userId = token.split(" ")[1];
        console.log('User ID extracted:', userId);  // Debugging step

        localStorage.setItem('token', token)
        localStorage.setItem('userId', userId)
        localStorage.setItem('username', this.username)

        this.$emit('login-success', {
          username: this.username,
          userId: userId,
          token: token,
        });

      } catch (error) {
        console.error('Errore di login:', error);
      }
    }
  }
};
</script>

<template>
  <div class="loginform">
    <img src="/wasaphoto.svg">
    <form @submit.prevent="dologin">
      <div>
        <input type="text" id="username" placeholder="Username" v-model="username">
      </div>
      <button type="submit">Login</button>
    </form>
  </div>
</template>

<style>
.loginform {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 90vh;
}

form {
  margin-top: 1rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

input {
  padding: 0.2rem;
  font-size: 1rem;
}

form button {
  padding: 0.2rem 1rem;
  font-size: 1rem;
  cursor: pointer;
  border-radius: 10px;
}

.loginform img{
  width: 300px;
  padding: 1rem;
}
</style>
