import http from "../http-common";

class BuyersDataService {
  getAll() {
    return http.get("/buyers");
  }
  getBuyerDetailedInformation(id) {
    return http.get("/buyers/"+id);
  }
}

export default new BuyersDataService();