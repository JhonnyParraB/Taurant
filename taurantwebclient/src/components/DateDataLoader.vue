<template>
  <v-row align="center" class="list px-3 mx-auto">
    <v-form v-model="valid">
        <div id="getUpdated">
            <h3>Get updated with the last information of your restaurant!</h3>
        </div>
        <div id="datePicker">
            <v-row justify="center">
                <v-date-picker color="#4b0082" v-model="picker"></v-date-picker>
            </v-row>
        </div>
        <div id="emailField">
        <v-text-field
            v-model="email"
            :rules="emailRules"
            label="E-mail"
            required
          ></v-text-field>
        </div>

        <div class="text-center">
    <v-dialog
      v-model="dialog"
      width="500"
    >
      <template v-slot:activator="{  }">
        <v-btn
            class="mr-4"
            :disabled="!valid"
            outlined
            color="indigo"
            @click="loadDateData()"
        >
          Get Updated
        </v-btn>
      </template>

      <v-card>
        <v-card-title class="headline primary">
          Request received
        </v-card-title>

        <v-card-text>
          Your request was received, we will send you a message to "{{this.email}}" with the result of the operation.
        </v-card-text>

        <v-divider></v-divider>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            color="primary"
            text
            @click="dialog = false"
          >
            Got it
          </v-btn>
        </v-card-actions>
            </v-card>
            </v-dialog>
        </div>
    </v-form>
  </v-row>
</template>


<script>
export default {
  name: "date-data-loader",
  data() {
    return {
      valid: false,
      dialog: false,
      email: '',
      picker: new Date().toISOString().substr(0, 10),
      emailRules: [
        v => !!v || 'E-mail is required',
        v => /.+@.+/.test(v) || 'E-mail must be valid',
      ],
    };
  },
  methods: {
    loadDateData(){
        this.dialog = true;
    },
  },
  mounted() {
    
  }
};
</script>

<style>
#getUpdated{
    margin-bottom: 5vh;
}
#emailField{
    margin-bottom: 5vh;
}
#datePicker{
    margin-bottom: 5vh;
}
#getUpdatedBtn{
    text-align: center;
}
</style>