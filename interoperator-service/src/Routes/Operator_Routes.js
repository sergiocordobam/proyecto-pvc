const express = require('express');
const OperatorController = require('../Controller/Operator_Controller');
const router = express.Router();

// router.get('/comunication/operators', OperatorController.fetchOperators.bind(OperatorController));
router.get('/operators', OperatorController.fetchOperators.bind(OperatorController));

module.exports = router;