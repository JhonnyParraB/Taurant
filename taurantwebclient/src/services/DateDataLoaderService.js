import http from "../http-common";

class DateDataLoaderService {
  loadDayData(time, email) {
    return http.post("/load-day-data/"+time+"?email="+email);
  }
}

export default new DateDataLoaderService();