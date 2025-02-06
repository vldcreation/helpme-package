const { GoogleGenerativeAI } = require('@google/generative-ai');

class GeminiClient {
  constructor(config) {
    if (!config.apiKey) {
      throw new Error('Gemini API key is required');
    }
    this.genAI = new GoogleGenerativeAI(config.apiKey);
  }

  async chat(message, options = {}) {
    try {
      const model = this.genAI.getGenerativeModel({ 
        model: options.model || 'gemini-1.5-pro'
      });

      const chat = model.startChat();
      const result = await chat.sendMessage(message);
      const response = await result.response;
      
      return response.text();
    } catch (error) {
      throw new Error(`GeminiClient chat error: ${error.message}`);
    }
  }
}

module.exports = GeminiClient;