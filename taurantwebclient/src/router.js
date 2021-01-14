import Vue from "vue";
import Router from "vue-router";

Vue.use(Router);

export default new Router({
  mode: "history",
  routes: [
    {
      path: "/buyers/:id",
      name: "buyer-information",
      component: () => import("./components/BuyerInformation")
    },
    {
      path: "/buyers",
      name: "buyers",
      component: () => import("./components/Buyers")
    },
    {
      path: "/",
      name: "home",
      component: () => import("./components/Buyers")
    },
  ]
});