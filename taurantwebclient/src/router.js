import Vue from "vue";
import Router from "vue-router";

Vue.use(Router);

export default new Router({
  mode: "history",
  routes: [
    {
      path: "/buyer/:id",
      alias: "/buyer",
      name: "buyer-information",
      component: () => import("./components/BuyerInformation")
    },
    {
      path: "/",
      alias: "/buyers",
      name: "buyers",
      component: () => import("./components/Buyers")
    },
  ]
});