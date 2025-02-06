const OpenAI = require('openai');

class OpenAIClient {
  constructor(config) {
    if (!config.apiKey) {
      throw new Error('OpenAI API key is required');
    }
    this.client = new OpenAI({
      apiKey: config.apiKey
    });
  }

  async chat(message, options = {}) {
    try {
      const completion = await this.client.chat.completions.create({
        model: options.model || 'gpt-3.5-turbo',
        messages: [{ role: 'user', content: message }],
        temperature: options.temperature || 0.7,
      });

      return completion.choices[0].message.content;
    } catch (error) {
      throw new Error(`OpenAIClient chat error: ${error.message}`);
    }
  }
}

module.exports = OpenAIClient;