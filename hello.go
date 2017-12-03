package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

//constantes
const monitoramentos = 3
const delay = 2

func main() {

	exibeIntroducao()
	//fmt.Println("variavel versão é do tipo", reflect.TypeOf(versao))

	//utilizamos um for sem passar nada para que o programa fique em loop e so saia quando escolhermos o 0
	for {
		exibeMenu()

		//utilizando o definidor de variavel curto é possivel atribuir funcoes
		comando := leComando()
		/*
			if em go deve sempre tornar um boleano e o else
			deve sempre estar grudando com a ultima chave

				if comando == 1 {
					fmt.Println("Monitorando...")
				} else if comando == 2 {
					fmt.Println("Exibindo Logs...")
				} else if comando == 0 {
					fmt.Println("Saindo do programa...")
				} else {
					fmt.Println("Não conheço este comando")
				}
		*/

		/*
		 deafult utilizado no switch como nas outras linguagens
		 caso não selecione nenhum dos casos é a vez dele,
		 e no Go o break é opciional para parar o comando
		 já que o 2º case só será executado se o 1º não for
		*/
		switch comando {
		case 1:
			iniciarMonitoramento()
		case 2:
			fmt.Println("Exibindo Logs...")
			imprimeLogs()
		case 0:
			fmt.Println("Saindo do programa...")

			/*
				através do pacote os é possivel comunicar com o sistema operacional
				passando uma mensagem 0 informando que quer encerrar o programa
			*/
			os.Exit(0)
		default:
			fmt.Println("Não conheço este comando")

			//passamos -1 para informar que houve um problema na execucao
			os.Exit(-1)
		}
	}

}
func exibeIntroducao() {
	/*
	   - scanF pede como parametro o marcador da variavel e um limitador
	   para saber oque se espera receber para aquela variavel, exemplo %d, para apenas numeros
	   - uma variavel deve ser declarada com :
	   var <nome da variavel> <tipo(opcional,se já atribuir um valor)>
	   porem podemos usar o atribuidor de variaveis curto:
	    <nome de variavel> := <valor a ser atribuido>
	*/
	nome := "Douglas"
	idade := 24
	versao := 1.1
	fmt.Println("Ola sr.", nome, "sua idade é", idade)
	fmt.Println("Este programa está na versão", versao)
}

func exibeMenu() {
	fmt.Println("1 - iniciar monitoramento")
	fmt.Println("2 - Exibir logs")
	fmt.Println("0 - Sair do programa")
}

func leComando() int {
	var comandoLido int

	/*
	   - scanF pede como parametro o marcador da variavel e um limitador
	   para saber oque se espera receber para aquela variavel, exemplo %d, para apenas numeros:
	   exemplo:
	   fmt.Scanf("%d", &comando)
	*/
	fmt.Scan(&comandoLido)
	fmt.Println("O comando escolhido foi", comandoLido)
	return comandoLido
}
func iniciarMonitoramento() {
	//através do pacote http dentro do pacote net é possivel fazer requsições http
	fmt.Println("Monitorando...")

	/*
			ao declarar um array devemos passar entre colchetes o numero de arrays e logo em seguida o tipo
			e se atentar pois não deve haver espaço entre ambos

		var sites [4]string
		sites[0] = "https://random-status-code.herokuapp.com/"
		sites[1] = "https://www.alura.com.br"
		sites[2] = "https://www.caelum.com.br"
		OBS: é possivel retornar um array através de um função
	*/

	/*
		-slices são como os arrays, mas tem tamanho dinamico,exemplo:
		nomes := []string
		-e é possivel atribuir diretamente:
		nomes := []string{"Douglas", "Daniel", "Bernardo"}
		-e conseguimos verificar seu tamanho através da função lens
		fmt.Println(len(nomes))
		-e podemos adicionar mais itens através da função append:
		nomes = append(nomes, "Aparecida")
		-alem disso podemos verificar a capacidade do array/slice através da funcao cap
		fmt.Println(cap(nomes))
	*/
	sites := leSitesDoArquivo()
	for i := 0; i < monitoramentos; i++ {
		//o range serve para iterar o nosso array e a cada iteração devolve a posição dele
		for i, site := range sites {
			fmt.Println("Testando site", i, ":", site)
			testaSite(site)
			fmt.Println()
		}
		//através do Sleep podemos pedir para que aconteça uma pause de tanto tempo até a proxima linha
		time.Sleep(delay * time.Minute)
	}
	fmt.Println("")
}

func testaSite(site string) {
	//o metodo GET retorna varios valores de uma requisição get
	//alem da resposta ele retorna um possivel erro
	resp, err := http.Get(site)
	if err != nil {
		fmt.Println("ocorreu o seguinte erro:", err)
	}
	if resp.StatusCode == 200 {
		fmt.Println("O Site", site, "foi carregado com sucesso")
		registraLogs(site, true)
	} else {
		fmt.Println("O Site", site, "está com problemas. Status Code", resp.StatusCode)
		registraLogs(site, false)
	}
}

func leSitesDoArquivo() []string {
	var sites []string

	//utilizamos os.Open para abrir arquivos em go
	// e o pacote bufio
	arquivo, err := os.Open("sites.txt")

	//tratamento de exceções em go
	if err != nil {
		fmt.Println("Ocorreu um erro", err)
	}
	//Bufio nos concede varias funcoes para leitura de arquivos
	leitor := bufio.NewReader(arquivo)
	/*
		dentre elas a ReadString que le uma linha do arquivo e retorna em String,
		 porem precisa de um delimitador, nesse caso o \n
	*/
	for {
		linha, err := leitor.ReadString('\n')
		//utilizamos o trimSpace do pacote strings para retirar os espaços
		linha = strings.TrimSpace(linha)
		sites = append(sites, linha)
		if err == io.EOF {
			break
		}
	}
	arquivo.Close()
	return sites
}
func registraLogs(site string, status bool) {
	/*
		para abrir um arquivo em go utilizamos a func OpenFile do pacote os
		e passamos algumas flags para dizer como o arquivo sera configurado
		flags podem ser observadas aqui
		https://golang.org/pkg/os/#pkg-constants
		através das flags usadas estamos dizendo para criar o arquivo se não existir,
		case exista escreva nele e se tiver algo escrito concatene
	*/
	arquivo, err := os.OpenFile("log.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
	}
	//conseguimos escrever no arquivo com a função WriteString
	//utilizamos a função FormatBool do pacote strconv, para formatar o nosso boleano em string
	//o seguinte trecho de codigo time.Now().Format("01/02/2006 03:04:05") é utilizado para formatar a data atual em go
	//OBS:a formação utiliza constantes no codigo para descobrir qual o tipo de formação vai ser, caso utilize um numero diferente resultará na não formataçao da data
	arquivo.WriteString(time.Now().Format("01/02/2006 03:04:05") + " - " + site + "- online: " + strconv.FormatBool(status) + "\n")
	arquivo.Close()
}

func imprimeLogs() {
	// utilizamos a função ReadFile do pacote ioUtil para ler um arquivo inteiro
	//OBS:nesse caso nao é necessario fechar o arquivo já que o pacote ioutil abre e fecha o arquivo para nós
	arquivo, err := ioutil.ReadFile("log.txt")
	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
	}
	fmt.Println(string(arquivo))
}

/*
exemplo de função que devolve mais de um valor:

func devolveNomeEIdade() (string, int) {
    nome := "Douglas"
    idade := 24
    return nome, idade
}

 e para receber os valores dessa função devemos utilizar duas variaveis:
 nome,idade := devolveNomeEIdade()

 quando utilizamos _, significa que queremos ignorar um desses retornos
 _,idade := devolveNomeEIdade()

*/
