import GeminiClient from './gemini/client.js';
import OpenAIClient from './openai/client.js';
import LocalClient from './local/client.js';
import ChatClientConfig from './config.js';
import {RESET, RED} from './error.js';

class ChatClient {
  constructor(agent='gemini', config={}, model='gemini-1.5-flash') {
    this.agent = agent;
    this.model = model;
    this.config = config;

    this.validateConfig();
    this.initializeClient();
  }

  validateConfig() {
    if (this.agent === 'gemini' && !this.config.apiKey) {
      throw new Error('Gemini API key is required');
    }
    if (this.agent === 'openai' && !this.config.apiKey) {
      throw new Error('OpenAI API key is required');
    }
  }

  initializeClient() {
    if (this.agent === 'gemini') {
      this.client = new GeminiClient(this.config);
    } else if (this.agent === 'openai') {
      this.client = new OpenAIClient(this.config);
    } else if (this.agent === 'local') {
      this.client = new LocalClient(this.config);
    } else {
      throw new Error(`Unsupported agent: ${this.agent}`);
    }
  }

  async chat(message, options = {}) {
    try {
      if (!this.client) {
        throw new Error('Client not initialized');
      }

      const defaultOptions = {
        model: this.model,
        ...options
      };

      return await this.client.chat(message, defaultOptions);
    } catch (error) {
      console.error(`Error in chat: ${error.message}`);
      throw error;
    }
  }

  // Method to switch agent at runtime
  switchAgent(newAgent, config = {}, model = null) {
    this.agent = newAgent;
    this.config = config;
    if (model) this.model = model;
    
    this.validateConfig();
    this.initializeClient();
  }
}

// Initialize with Gemini
const config = new ChatClientConfig();
const chatClient = new ChatClient('openai', {
  apiKey: config.gemini.apiKey,
});

// validate
try {
  config.validate();
}
catch (error) {
 console.error(`${RED}Error: ${error.message}${RESET}`);
}

// Chat using Gemini
console.log(await chatClient.chat("How are you?"));

// Initialize with Local
// const config = new ChatClientConfig();
// const chatClient = new ChatClient('local', config.local);

// Switch to local agent
// chatClient.switchAgent('local', {
//   url: config.local.url,
//   model: config.local.model
// });

// Chat using local agent
// chatClient.chat("Now using local agent");