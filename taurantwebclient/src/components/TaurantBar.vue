<template>
  <div>
    <button @click="goHome()" id="taurantBtn">
         Taurant 
    </button>
  </div>
</template>


<script>
export default {
  name: "taurant-bar",
  data() {
  },
  methods: {
    goHome(){
      this.$router.push({name: "buyers"}).catch(()=>{
        this.$router.go()
      });
    }
  },
  mounted() {
    
  }
};
</script>

<style>
  @import url('https://fonts.googleapis.com/css2?family=Pacifico&display=swap');

  #taurantBtn {
    border: none;
    padding: 0;
    background: none;
    outline:none;
    font-family: 'Pacifico', cursive;
    font-size: 6vh;
    user-select: none;
    margin-left: 3vw;
  }
</style>
