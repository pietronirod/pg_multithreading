# 🚀 Desafio de Multithreading e APIs em Go

Este repositório apresenta uma solução robusta e eficiente para um desafio técnico que explora as capacidades de multithreading e integração com APIs utilizando Go. O objetivo é determinar, entre duas APIs concorrentes, qual delas retorna os dados mais rapidamente e descartar a resposta mais lenta.

## 📋 Descrição do Desafio

Neste desafio, você precisará utilizar multithreading para realizar requisições simultâneas a duas APIs distintas e escolher a resposta mais rápida:

- **BrasilAPI:** <https://brasilapi.com.br/api/cep/v1/> + cep
- **ViaCEP:** <http://viacep.com.br/ws/"> + cep + /json/

### 🎯 Requisitos

1. **Priorizar a Velocidade:** A API que entregar a resposta mais rápida deve ser acatada, enquanto a resposta mais lenta deve ser descartada.
2. **Exibição dos Resultados:** O resultado deve ser exibido no terminal, incluindo os dados do endereço e a fonte da API utilizada.
3. **Limite de Tempo:** Se nenhuma das APIs responder em até 1 segundo, um erro de timeout deve ser exibido.

## ⚙️ Execução do Projeto

### Pré-requisitos

Para executar esta aplicação, você precisará ter o Go (versão 1.16 ou superior) instalado em seu ambiente de desenvolvimento.

### Passos para Execução

1. **Clone o Repositório**

   ```bash
   git clone https://github.com/seu_usuario/seu_repositorio.git
   cd seu_repositorio
   ```

2. **Execute a Aplicação** Execute o código para ver qual API responde mais rápido:

    ```bash
    go run main.go
    ```

3. **Execute os Testes** Verifique a robustez da implementação através de testes unitários:

    ```bash  
    go test -v
    ```

## 🛠️ Configuração Personalizada

Este projeto permite configurações flexíveis através de variáveis de ambiente:

- `BRASIL_API_URL`: Defina a URL da BrasilAPI (padrão: <https://brasilapi.com.br/api/cep/v1/>).
- `VIACEP_URL`: Defina a URL da ViaCEP (padrão: <https://viacep.com.br/ws/>).
- `API_TIMEOUT`: Defina o tempo limite para as requisições (padrão: 1s).

Essas configurações permitem ajustar o comportamento da aplicação para diferentes ambientes e necessidades.

## 🧩 Considerações Técnicas

Este projeto foi desenvolvido com foco em alta performance e resiliência, utilizando conceitos avançados de programação concorrente em Go. A solução demonstra como o uso eficiente de goroutines pode otimizar a latência de sistemas que dependem de múltiplos serviços externos.

Além disso, a implementação inclui tratamento de erros detalhado, permitindo uma depuração mais eficaz e um melhor entendimento do comportamento da aplicação sob diferentes condições.

## 🤝 Contribuições

Contribuições são bem-vindas! Se você tiver sugestões ou melhorias, sinta-se à vontade para abrir uma issue ou enviar um pull request.

---
Desenvolvido com Go para garantir desempenho, eficiência e simplicidade. Vamos construir algo incrível juntos! 🚀
