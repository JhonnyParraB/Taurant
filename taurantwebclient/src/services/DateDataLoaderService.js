import http from "../http-common";

class DateDataLoaderService {
  loadDayData(time) {
    return http.get("/load-day-data/"+time);
  }
}

export default new DateDataLoaderService();