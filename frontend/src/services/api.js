import axios from 'axios';

const api = axios.create({
  baseURL: process.env.REACT_APP_API_URL || 'http://localhost:8080/api',
  headers: {
    'Content-Type': 'application/json',
  },
});

const apiService = {
  checkVoucherExists: async (flightNumber, date) => {
    try {
      const response = await api.post('/check', {
        flightNumber,
        date
      });
      return response.data;
    } catch (error) {
      throw new Error(error.response?.data?.error || 'Failed to check voucher existence');
    }
  },

  generateVoucher: async (voucherData) => {
    try {
      const response = await api.post('/generate', voucherData);
      return response.data;
    } catch (error) {
      throw new Error(error.response?.data?.error || 'Failed to generate vouchers');
    }
  },

  healthCheck: async () => {
    try {
      const response = await api.get('/health');
      return response.data;
    } catch (error) {
      throw new Error('Health check failed');
    }
  }
};

export default apiService; 