const axios = require('axios');
require('dotenv').config({ path: '../config/dev.env' });
class OperatorService {
    constructor() {
        this.api_url = process.env.API_BASE_URL;
        this.operator_name = process.env.OPERATOR_NAME;
        this.operator_address = process.env.OPERATOR_ADDRESS;
        this.operator_participants = process.env.OPERATOR_PARTICIPANTS.split(',').map(participant => participant.trim());
        this.contactMail = process.env.OPERATOR_CONTACT_EMAIL;
    }
    async getOperators() {
        try {
            const response = await axios.get(`${this.api_url}/getOperators`);
            return response.data;
        } catch (error) {
            console.error('Error fetching data from API:', error.message);
            throw error;
        }
    }

    async getOperatorByName(name) {
        const response = await this.getOperators();
        const operator = response.find((op) => op.name === name);
        return operator
    }

    async getOperatorById(id) {
        const response = await this.getOperators();
        const operator = response.find((op) => op.id === id);
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
            return response.data;
        } catch (error) {
            console.error('Error registering operator:', error.message);
            throw error;
        }
    }

    async checkOperatorInSystem(){
        const operator = await getOperatorByName(this.operator_name);
        if (operator){
            this.saveToken(operator.OperatorId);
        }
        else{
            const register = this.registerOperator();
            this.saveToken(register);
        }
    }

    async saveToken(token) {
        process.env.OPERATOR_ID = token;
    }
}

module.exports = OperatorService;