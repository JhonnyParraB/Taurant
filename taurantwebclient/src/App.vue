<template>
  <v-app>
  <div id="MainDiv">
    <div class="main">
      <div class="top">
        <taurant-bar></taurant-bar>
      </div>
      <div id="lower">
        <div class="left">
          <date-data-loader></date-data-loader>
        </div>
        <div class="right">
          <router-view></router-view>
        </div>
      </div>
    </div>
  </div>
  </v-app>
</template>

<script>
import DateDataLoader from './components/DateDataLoader.vue';
import TaurantBar from './components/TaurantBar.vue';

export default {
  components: { DateDataLoader, TaurantBar },
  name: 'App',
};
</script>

<style scoped>
#MainDiv {
  width: 100vw;
  height: 100vh;
  overflow:auto;
}

.main {
  width: auto;
  display: flex;
  flex-direction: column;
  overflow:auto;
}

#lower {
  display: flex;
  flex-direction: row;
  overflow:auto;
}

.top {
  width: 100vw;
  height: 10vh;
  background: rgb(0,0,0);
  background: linear-gradient(90deg, rgba(0,0,0,1) 5%, rgba(75,1,130,1) 48%);
  border-bottom: 1px solid gray;
  padding: 10px;
  position: absolute;
}

.left{
  width: 30vw;
  height: 90vh;
  padding: 10px;
  padding-top: 10vh;
  margin-top: 10vh;
  border-right: 2px solid gray;
  position:absolute;
}
.right {
  width: 70vw;
  height: 90vh;
  margin-left: 30vw;
  margin-top: 10vh;
  padding: 10px;
  overflow-y: scroll;
}
</style>
