const axios = require('axios');
require('dotenv').config({ path: './Config/dev.env' });
class OperatorService {
    constructor() {
        this.api_url = process.env.API_BASE_URL;
        this.operator_name = process.env.OPERATOR_NAME;
        this.operator_address = process.env.OPERATOR_ADDRESS;
        this.operator_participants = [process.env.OPERATOR_PARTICIPANTS1,
                                    process.env.OPERATOR_PARTICIPANTS2, process.env.OPERATOR_PARTICIPANTS3
        ];
        this.contactMail = process.env.OPERATOR_CONTACT_EMAIL;
    }
    async getOperators() {
        try {
            console.log(`Fetching operators from URL: ${this.api_url}/getOperators`);
            const response = await axios.get(`${this.api_url}/getOperators`);
            
            console.log("Operator service initialized and all operators fetched.");
            return response.data;
        } catch (error) {
            console.error('Error fetching data from API:', error.message);
            throw error;
        }
    }

    async getOperatorByName(name) {
        const response = await this.getOperators();
        const operator = response.find((op) => op.name === name);
        console.log(`${name} Operator fetched by operator service`);
        return operator
    }

    async getOperatorById(id) {
        const response = await this.getOperators();
        const operator = response.find((op) => op.id === id);
        console.log(`${id} Operator fetched by operator service`);
        return operator
    }

    async registerOperator() {
        const operator = {
            name: this.operator_name,
            address: this.operator_address,
            contactMail: this.contactMail,
            participants: this.operator_participants
        };
        try {
            const response = await axios.post(`${this.api_url}/registerOperator`, operator);
            console.log(`Operator registered successfully: ${this.operator_name}`);
            return response.data;
        } catch (error) {
            console.error('Error registering operator:', error.message);
            throw error;
        }
    }

    async checkOperatorInSystem(){
        const operator = await this.getOperatorByName(this.operator_name);
        try{
            if (operator){
                console.log(`Operator ${this.operator_name} already exists in the system.`);
                return operator.OperatorId;
            }
            else{
                const register = this.registerOperator();
                console.log(`Operator ${this.operator_name} does not exist in the system. Registering...`);
                this.saveToken(register);
                console.log(`Operator ${this.operator_name} registered successfully and ID saved.`);
            }
        }catch (error) {
            console.error('Error checking operator in system:', error.message);
            throw error;
        }
    }

    async saveToken(token) {
        process.env.OPERATOR_ID = token;
    }
}

module.exports = OperatorService;