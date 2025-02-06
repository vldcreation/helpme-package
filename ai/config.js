class ChatClientConfig {
  constructor() {
    this.gemini = {
      apiKey: process.env.GEMINI_API_KEY || '',
      model: process.env.GEMINI_MODEL || 'gemini-1.5-pro'
    };

    this.openai = {
      apiKey: process.env.OPENAI_API_KEY || '',
      model: process.env.OPENAI_MODEL || 'gpt-3.5-turbo'
    };

    this.local = {
      url: process.env.LOCAL_AI_URL || 'http://localhost:6969',
      model: process.env.LOCAL_AI_MODEL || 'gemini-1.5-flash'
    };
  }

  validate() {
    if (!this.gemini.apiKey && !this.openai.apiKey && !this.local.url) {
      throw new Error('At least one AI service configuration must be provided');
    }
  }
}

module.exports = ChatClientConfig;