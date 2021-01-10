<template>
  <v-row align="center" class="list px-3 mx-auto">
    <v-col cols="12" sm="12">
      <v-card class="mx-auto" tile>
        <v-card-title>Buyer Information</v-card-title>
        <v-card>
            <v-card-title>{{ name }}</v-card-title>
            <v-card-subtitle>{{ age }} years</v-card-subtitle>
        </v-card>
        <v-data-table
          :headers="headers"
          :items="buyers"
          disable-pagination
          :hide-default-footer="true"
          @click:row="seeBuyerInformation"
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
      buyerInformation: null,
      headers: [
        { text: "ID", align: "start", sortable: true, value: "id" },
        { text: "Name", sortable: true, value: "name" },
        { text: "Age", value: "age", sortable: true },
      ],
    };
  },
  methods: {
    retrieveBuyerInformation() {
      BuyersDataService.getAll()
        .then((response) => {
          this.buyerInformation = response
        })
        .catch((e) => {
          console.log(e);
        });
    },
    refreshList() {
      this.retrieveBuyerInformation();
    },
  },
  mounted() {
    this.retrieveBuyers();
  }
};
</script>

<style>
.list {
  max-width: 750px;
}
</style>