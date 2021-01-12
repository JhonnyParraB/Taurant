import Vue from "vue";
import Router from "vue-router";

Vue.use(Router);

export default new Router({
  mode: "history",
  routes: [
    {
      path: "/buyers/:id",
      alias: "/buyer",
      name: "buyer-information",
      component: () => import("./components/BuyerInformation")
    },
    {
      path: "/buyers",
      alias: "/buyers",
      name: "buyers",
      component: () => import("./components/Buyers")
    },
    {
      path: "/",
      alias: "/buyers",
      name: "buyers",
      component: () => import("./components/Buyers")
    },
  ]
});