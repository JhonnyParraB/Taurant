<template>
  <v-row align="center" class="list px-3 mx-auto">



    <v-col cols="12" sm="12">
      <v-card class="mx-auto" tile>
        <v-card-title>Buyers</v-card-title>

        <v-data-table
          :headers="headers"
          :items="buyers"
          disable-pagination
          :hide-default-footer="true"
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
      title: "",
      headers: [
        { text: "Name", align: "start", sortable: true, value: "name" },
        { text: "Age", value: "age", sortable: true },
      ],
    };
  },
  methods: {
    retrieveBuyers() {
      BuyersDataService.getAll()
        .then((response) => {
          this.buyers = response.data.map(this.getDisplayBuyer);
          console.log(response.data);
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
      };
    },
  },
  mounted() {
    this.retrieveBuyers();
  },
};
</script>

<style>
.list {
  max-width: 750px;
}
</style>1