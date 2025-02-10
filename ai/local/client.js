class LocalClient {
  constructor(config) {
    this.url = config.url || 'http://localhost:6969';
    this.model = config.model || 'gemini-1.5-flash';
  }

  async chat(message, options = {}) {
    try {
      const response = await fetch(`${this.url}/v1/chat/completions`, {
        method: 'POST',
        headers: {
            'accept': 'application/json',
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          model: options.model || this.model,
          messages: [{ role: 'user', content: message }]
        })
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
    //   console.log('LocalClient chat:', data.choices[0].message.content)
      return data.choices[0].message.content;
    } catch (error) {
      throw new Error(`LocalClient chat error: ${error.message}`);
    }
  }
}

export default LocalClient;