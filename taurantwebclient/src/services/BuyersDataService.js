import http from "../http-common";

class BuyersDataService {
  getAll() {
    return http.get("/buyers");
  }
}

export default new BuyersDataService();