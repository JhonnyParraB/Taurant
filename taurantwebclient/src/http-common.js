import axios from "axios";

export default axios.create({
  baseURL: process.env.VUE_APP_TAURANT_SERVER_URL,
  headers: {
    "Content-type": "application/json"
  }
});