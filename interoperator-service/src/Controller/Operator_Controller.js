const OperatorService = require("../Service/Operator_Service");

class OperatorController {
    constructor() {
        this.operatorService = new OperatorService();
        this.initialize();
    }

    async initialize() {
        try {
            await this.operatorService.checkOperatorInSystem();
            console.log("Operator service initialized and all operators fetched.");

        } catch (error) {
            console.error("Error initializing operator service:", error.message);
        }
    }

    async fetchOperators(req, res) {
        try {
            const operators = await this.operatorService.getOperators();
            res.status(200).json(operators);
        } catch (error) {
            console.error('Error fetching operators:', error.message);
            res.status(500).json({ error: 'Failed to fetch operators' });
        }
    }    
}

module.exports = new OperatorController();
