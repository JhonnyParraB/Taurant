<template>
  <v-row align="center" class="list px-3 mx-auto">
    <v-col cols="12" sm="12">
      <v-card class="mx-auto" tile>
        <v-card-title class="justify-center">Buyers</v-card-title>
        <v-data-table
          :headers="headers"
          :items="buyers"
          disable-pagination
          :hide-default-footer="true"
          @click:row="seeBuyerInformation"
          :loading="loading"
        >
        </v-data-table>
      </v-card>
    </v-col>
  </v-row>
</template>

<script>
import BuyersDataService from "../services/BuyersDataService";
export default {
  name: "buyers-list",
  data() {
    return {
      buyers: [],
      loading: true,
      title: "",
      headers: [
        { text: "ID", align: "start", sortable: true, value: "id" },
        { text: "Name", sortable: true, value: "name" },
        { text: "Age", value: "age", sortable: true },
      ],
    };
  },
  methods: {
    retrieveBuyers() {
      BuyersDataService.getAll()
        .then((response) => {
          this.loading = false;
          this.buyers = response.data.map(this.getDisplayBuyer);
        })
        .catch((e) => {
          console.log(e);
        });
    },
    refreshList() {
      this.retrieveBuyers();
    },
    getDisplayBuyer(buyer) {
      return {
        name: buyer.name,
        age: buyer.age,
        id: buyer.id
      };
    },
    seeBuyerInformation(item){
      this.$router.push({name: "buyer-information",  params: {id:item.id}});
    }
  },
  mounted() {
    console.log(process.env.VUE_APP_TAURANT_SERVER_URL)
    this.retrieveBuyers();
  }
};
</script>

<style>
.list {
  max-width: 750px;
}
</style>