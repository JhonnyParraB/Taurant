<template>
  <v-row align="center" class="list px-3 mx-auto">
    <v-col cols="12" sm="12">
      <v-card class="mx-auto" tile>
        <v-card-title>Buyer Information</v-card-title>
        <v-card>
            <v-card-title>{{ this.buyerInformation.buyer.name }}</v-card-title>
            <v-card-subtitle>{{ this.buyerInformation.buyer.age }} years</v-card-subtitle>

            <v-expansion-panels popout>
              <v-expansion-panel>
                <v-expansion-panel-header>Shopping History</v-expansion-panel-header>
                <v-expansion-panel-content>
                  <v-data-table
                    :headers="transactionsHeaders"
                    :items="this.buyerInformation.buyer.transactions"
                    disable-pagination
                    :hide-default-footer="true"
                  >
                  </v-data-table>
                </v-expansion-panel-content>
              </v-expansion-panel>

              <v-expansion-panel>
                <v-expansion-panel-header>Buyers Using the Same IP</v-expansion-panel-header>
                <v-expansion-panel-content>
                  <v-data-table
                    :headers="buyersWithSameIPHeaders"
                    :items="this.buyerInformation.buyers_with_same_ip"
                    disable-pagination
                    :hide-default-footer="true"
                  >
                  </v-data-table>
                </v-expansion-panel-content>
              </v-expansion-panel>
              <v-expansion-panel>
                <v-expansion-panel-header>Recommended Products</v-expansion-panel-header>
                <v-expansion-panel-content>
                  <v-data-table
                    :headers="recommendedProductsHeader"
                    :items="this.buyerInformation.recommended_products"
                    disable-pagination
                    :hide-default-footer="true"
                  >
                  </v-data-table>
                </v-expansion-panel-content>
              </v-expansion-panel>
            </v-expansion-panels>
        </v-card>
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
      buyer_id: "",
      transactionsHeaders: [
        { text: "ID", align: "start", sortable: true, value: "id" },
        { text: "IP", sortable: true, value: "location.ip" },
        { text: "Device", value: "device", sortable: true },
      ],
      buyersWithSameIPHeaders: [
        { text: "ID", align: "start", sortable: true, value: "id" },
        { text: "Name", sortable: true, value: "name" },
        { text: "Age", value: "age", sortable: true },
        { text: "Shared IP", value: "shared_ip", sortable: true },
      ],
      recommendedProductsHeader: [
        { text: "ID", align: "start", sortable: true, value: "id" },
        { text: "Name", sortable: true, value: "name" },
        { text: "Price", value: "price", sortable: true },
      ],
    };
  },
  methods: {
    retrieveBuyerInformation() {
      this.buyer_id = this.$route.params.id
      BuyersDataService.getBuyerDetailedInformation(this.buyer_id)
        .then((response) => {
          this.buyerInformation = response.data;
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
    this.retrieveBuyerInformation();
  }
};
</script>

<style>
.list {
  max-width: 750px;
}
</style>