import React, { useState } from 'react';
import apiService from '../services/api';

const VoucherGenerator = () => {
  const [formData, setFormData] = useState({
    name: '',
    id: '',
    flightNumber: '',
    date: '',
    aircraft: 'ATR'
  });
  
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [seats, setSeats] = useState([]);
  const [showResults, setShowResults] = useState(false);
  const [regenerateMode, setRegenerateMode] = useState(false);
  const [selectedSeats, setSelectedSeats] = useState([]);

  const aircraftTypes = ['ATR', 'Airbus 320', 'Boeing 737 Max'];

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const formatDateForDisplay = (dateString) => {
    if (!dateString) return '';
    const date = new Date(dateString);
    const day = String(date.getDate()).padStart(2, '0');
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const year = String(date.getFullYear()).slice(-2);
    return `${day}-${month}-${year}`;
  };

  const formatDateForAPI = (dateString) => {
    return dateString;
  };

  const toggleSeatSelection = (seat) => {
    setSelectedSeats(prev => {
      if (prev.includes(seat)) {
        return prev.filter(s => s !== seat);
      } else {
        return [...prev, seat];
      }
    });
  };

  const handleRegenerate = async () => {
    if (selectedSeats.length === 0) {
      setError('Please select at least one seat to regenerate');
      return;
    }

    setLoading(true);
    setError('');

    try {
      const apiData = {
        ...formData,
        date: formatDateForAPI(formData.date)
      };

      const regenerateResponse = await apiService.regenerateVoucher(apiData, selectedSeats);
      
      if (regenerateResponse.success) {
        setSeats(regenerateResponse.seats);
        setRegenerateMode(false);
        setSelectedSeats([]);
      } else {
        setError('Failed to regenerate vouchers');
      }
      
    } catch (err) {
      setError(err.message || 'An error occurred while regenerating vouchers');
    } finally {
      setLoading(false);
    }
  };

  const resetToGenerate = () => {
    setFormData({
      name: '',
      id: '',
      flightNumber: '',
      date: '',
      aircraft: 'ATR'
    });
    setSeats([]);
    setShowResults(false);
    setRegenerateMode(false);
    setSelectedSeats([]);
    setError('');
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!formData.name || !formData.id || !formData.flightNumber || !formData.date || !formData.aircraft) {
      setError('Please fill in all required fields');
      return;
    }

    setLoading(true);
    setError('');
    setSeats([]);
    setShowResults(false);
    setRegenerateMode(false);
    setSelectedSeats([]);

    try {
      const apiData = {
        ...formData,
        date: formatDateForAPI(formData.date)
      };

      const checkResponse = await apiService.checkVoucherExists(apiData.flightNumber, apiData.date);
      
      if (checkResponse.exists) {
        setError('Vouchers have already been generated for this flight date');
        setLoading(false);
        return;
      }

      const generateResponse = await apiService.generateVoucher(apiData);
      
      if (generateResponse.success) {
        setSeats(generateResponse.seats);
        setShowResults(true);
      } else {
        setError('Failed to generate vouchers');
      }
      
    } catch (err) {
      setError(err.message || 'An error occurred while generating vouchers');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-100 py-8 px-4">
      <div className="max-w-2xl mx-auto">
        <div className="bg-white rounded-lg shadow-md p-6">
          <h1 className="text-2xl font-bold text-gray-800 mb-6">
            Crew Voucher Generator
          </h1>
          
          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Crew Name *
              </label>
              <input
                type="text"
                name="name"
                value={formData.name}
                onChange={handleChange}
                required
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="Enter crew name"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Crew ID *
              </label>
              <input
                type="text"
                name="id"
                value={formData.id}
                onChange={handleChange}
                required
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="Enter crew ID"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Flight Number *
              </label>
              <input
                type="text"
                name="flightNumber"
                value={formData.flightNumber}
                onChange={handleChange}
                required
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="Enter flight number"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Flight Date * (DD-MM-YY)
              </label>
              <input
                type="date"
                name="date"
                value={formData.date}
                onChange={handleChange}
                required
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
              {formData.date && (
                <p className="text-sm text-gray-600 mt-1">
                  Display format: {formatDateForDisplay(formData.date)}
                </p>
              )}
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Aircraft Type *
              </label>
              <select
                name="aircraft"
                value={formData.aircraft}
                onChange={handleChange}
                required
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                {aircraftTypes.map(type => (
                  <option key={type} value={type}>{type}</option>
                ))}
              </select>
            </div>

            {error && (
              <div className="bg-red-50 border border-red-200 rounded-md p-4">
                <div className="flex">
                  <div className="flex-shrink-0">
                    <svg className="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
                      <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
                    </svg>
                  </div>
                  <div className="ml-3">
                    <p className="text-sm font-medium text-red-800">{error}</p>
                  </div>
                </div>
              </div>
            )}

            <button
              type="submit"
              disabled={loading}
              className="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              {loading ? 'Generating Vouchers...' : 'Generate Vouchers'}
            </button>
          </form>

          {showResults && seats.length > 0 && (
            <div className="mt-8 p-6 bg-green-50 border border-green-200 rounded-md">
              <div className="flex items-center justify-between mb-4">
                <h2 className="text-lg font-semibold text-green-800">
                  Vouchers Generated Successfully!
                </h2>
                <div className="flex space-x-2">
                  {!regenerateMode && (
                    <button
                      onClick={() => setRegenerateMode(true)}
                      className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition-colors"
                    >
                      Regenerate Seats
                    </button>
                  )}
                  <button
                    onClick={resetToGenerate}
                    className="px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700 transition-colors"
                  >
                    New Voucher
                  </button>
                </div>
              </div>

              {regenerateMode && (
                <div className="mb-4 p-4 bg-yellow-50 border border-yellow-200 rounded">
                  <h3 className="font-medium text-yellow-800 mb-2">
                    Regenerate Mode: Click on seats to select them for regeneration
                  </h3>
                  <div className="flex space-x-2">
                    <button
                      onClick={handleRegenerate}
                      disabled={loading || selectedSeats.length === 0}
                      className="px-4 py-2 bg-orange-600 text-white rounded hover:bg-orange-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                    >
                      {loading ? 'Regenerating...' : `Regenerate ${selectedSeats.length} Seat(s)`}
                    </button>
                    <button
                      onClick={() => {
                        setRegenerateMode(false);
                        setSelectedSeats([]);
                      }}
                      className="px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700 transition-colors"
                    >
                      Cancel
                    </button>
                  </div>
                  {selectedSeats.length > 0 && (
                    <p className="text-sm text-yellow-700 mt-2">
                      Selected seats: {selectedSeats.join(', ')}
                    </p>
                  )}
                </div>
              )}

              <div className="space-y-2">
                {seats.map((seat, index) => (
                  <div 
                    key={index} 
                    className={`flex items-center justify-between bg-white p-3 rounded border ${
                      regenerateMode 
                        ? 'cursor-pointer hover:bg-gray-50 transition-colors' 
                        : ''
                    } ${
                      selectedSeats.includes(seat) 
                        ? 'ring-2 ring-orange-500 bg-orange-50' 
                        : ''
                    }`}
                    onClick={() => regenerateMode && toggleSeatSelection(seat)}
                  >
                    <div className="flex items-center">
                      <div className={`w-8 h-8 rounded-full flex items-center justify-center font-medium mr-3 ${
                        selectedSeats.includes(seat) 
                          ? 'bg-orange-100 text-orange-600' 
                          : 'bg-blue-100 text-blue-600'
                      }`}>
                        {regenerateMode && selectedSeats.includes(seat) ? '✓' : index + 1}
                      </div>
                      <div>
                        <p className="font-medium text-gray-800">Seat {seat}</p>
                        <p className="text-sm text-gray-600">Assigned to: {formData.id}</p>
                      </div>
                    </div>
                    <div className="text-sm font-medium text-green-600">
                      {regenerateMode ? (
                        <span className="text-orange-600">
                          {selectedSeats.includes(seat) ? 'Selected' : 'Click to select'}
                        </span>
                      ) : (
                        '✓ Voucher Created'
                      )}
                    </div>
                  </div>
                ))}
              </div>
              <div className="mt-4 p-3 bg-blue-50 border border-blue-200 rounded">
                <p className="text-sm text-blue-800">
                  <strong>Flight:</strong> {formData.flightNumber} | 
                  <strong> Date:</strong> {formatDateForDisplay(formData.date)} | 
                  <strong> Aircraft:</strong> {formData.aircraft}
                </p>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default VoucherGenerator; 