<template>
  <v-row align="center" class="list px-3 mx-auto">
    <v-col cols="12" sm="12">
      <v-card class="mx-auto" tile>
        <v-card-title class="justify-center">Buyer Information</v-card-title>
        <v-card>
            <v-card-title>{{ this.buyerInformation.buyer.id }}    {{this.buyerInformation.buyer.name}}</v-card-title>
            <v-card-subtitle>{{ this.buyerInformation.buyer.age }} years</v-card-subtitle>

            <v-expansion-panels>
              <v-expansion-panel>
                <v-expansion-panel-header>Shopping History</v-expansion-panel-header>
                <v-expansion-panel-content>
                  <v-data-table
                    :headers="transactionsHeaders"
                    :items="this.buyerInformation.buyer.transactions"
                    disable-pagination
                    :hide-default-footer="true"
                    show-expand
                  >
                  <template v-slot:expanded-item="{ headers, item }">
                    <td :colspan=headers.length>
                      <v-data-table
                      :headers="productOrderHeaders"
                      :items="item.product_orders"
                      disable-pagination
                      :hide-default-footer="true"
                      class="elevation-1 primary"
                      >
                        <template v-slot:[`item.subtotal`]="{ item }">
                          <span> {{ centsToDollars(parseInt(item.product.price) * parseInt(item.quantity)) }}</span>
                        </template>
                        <template v-slot:[`item.product.price`]="{ item }">
                          <span>{{ centsToDollars(item.product.price) }}</span>
                        </template>
                        <template slot="body.append">
                          <tr>
                              <th>Total</th>
                              <th></th>
                              <th></th> 
                              <th></th>
                              <th> {{ getTotalBuy(item.product_orders) }} </th>
                          </tr>                   
                        </template>
                      </v-data-table>
                    </td> 
                  </template>
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
                    @click:row="seeBuyerInformation"
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
                    <template v-slot:[`item.price`]="{ item }">
                      <span>{{ centsToDollars(item.price) }}</span>
                    </template>
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
  name: "buyer-information",
  data() {
    return {
      buyerInformation: null,
      buyer_id: "",
      transactionsHeaders: [
        { text: "ID", align: "start", sortable: true, value: "id" },
        { text: "IP", sortable: true, value: "location.ip" },
        { text: "Device", value: "device", sortable: true },
        { text: '', value: 'data-table-expand' },
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
        { text: "Price ($USD)", value: "price", sortable: true },
      ],
      productOrderHeaders: [
        { text: "ID", align: "start", sortable: true, value: "product.id" },
        { text: "Name", sortable: true, value: "product.name" },
        { text: "Price ($USD)", value: "product.price", sortable: true },
        { text: "Quantity", value: "quantity", sortable: true },
        { text: "Subtotal ($USD)", value: "subtotal", sortable: true },
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
    centsToDollars (value) {
        return parseInt(value)/100;
    },
    getTotalBuy (product_orders) {
        var total = 0;
        for(var i in product_orders){
          total += this.centsToDollars(parseInt(product_orders[i].product.price) * parseInt(product_orders[i].quantity)) 
        }
        return total
    },
    seeBuyerInformation(item){
      this.$router.push({name: "buyer-information",  params: {id:item.id}});
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