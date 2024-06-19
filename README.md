# Go Concurrency Patterns

O objetivo deste repositório é manter salvo os fontes de um artigo que estou escrevendo com alguns amigos.

## Artigo

Existem alguns padrões de utilização de concorrência em Go estabelecidos para resolver problemas comuns encontrados na programação concorrente.

### Worker Pools - Distribuição inteligente de tarefas

Worker pools permitem distribuir uma tarefa entre uma quantidade pré definida de workers que compartilham o mesmo estado (channels) de forma auto gerenciada.

```go
//Worker que chama a function "count" e escreve a resposta no channel "results"
func worker(tasks <-chan int, results chan<- int) {
 for task := range tasks {
  results <- count(task) 
 }
}

//fuction simples que conta até o número "n"
func count(n int) int {
 c := 1
 for i := 1; i <= n; i++ {
  c++
 }
 return c
}

func start(workersCount int, maxNumber int) (time.Duration, []int) {
 start := time.Now()
 tasks := make(chan int, maxNumber)
 results := make(chan int, maxNumber)
 nums := make([]int, maxNumber)

 //executa os workers
 for i := 0; i < workersCount; i++ {
  go worker(tasks, results)
 }

 //escrever no channel "tasks" vai disparar a contagem dos workers
 for i := 0; i < maxNumber; i++ {
  tasks <- i
 }

 //encerra os channels tasks
 close(tasks)

 //lê o que está chegando no channel "results" e popula o array
 for i := 0; i < maxNumber; i++ {
  results := <-results
  nums[i] = results
 }
 //calcula o tempo de execução
 elapsed := time.Since(start)
 return elapsed, nums
}

func main() {
 maxNumber := 10000
 timeOneWorker, numsOw := start(1, maxNumber)
 timeTwoWorkers, _ := start(2, maxNumber)
 timeThreeWorkers, _ := start(3, maxNumber)
 tomeFourWorkers, numsFw := start(4, maxNumber)

 fmt.Println("1 worker: ", timeOneWorker)
 fmt.Println("2 workers: ", timeTwoWorkers)
 fmt.Println("3 workers: ", timeThreeWorkers)
 fmt.Println("4 workers: ", tomeFourWorkers)

 fmt.Println("-----------------")
 fmt.Println("1 Workers", numsOw[:100])
 fmt.Println("4 Workers", numsFw[:100])

 fmt.Println("end")
}
```

Saída

```bash
go run worker-pool.go 
1 worker: 18.417556ms
2 workers: 9.794058ms
3 workers: 7.220087ms
4 workers: 6.473782ms
-----------------
1 Workers [1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31 32 33 34 35 36 37 38 39 40 41 42 43 44 45 46 47 48 49 50 51 52 53 54 55 56 57 58 59 60 61 62 63 64 65 66 67 68 69 70 71 72 73 74 75 76 77 78 79 80 81 82 83 84 85 86 87 88 89 90 91 92 93 94 95 96 97 98 99 100]
4 Workers [1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31 32 33 34 35 36 37 38 39 40 41 42 43 44 45 46 47 48 49 50 51 52 53 54 55 56 57 58 59 60 61 62 63 64 65 66 67 68 69 70 71 72 73 74 75 76 77 78 79 80 81 82 83 84 85 86 87 88 89 90 91 92 93 94 95 96 97 98 99 100]
end
```

### Pipeline - Processamento em estágios concorrentes

Pipeplines permitem o processamento de dados em estágios distintos, esses estágios podem rodar concorrentemente e sob demanda utilizando os poder das Go routines.

```go
// Cria um canal de inteiros e o preenche com os valores do slice de inteiros
func generate(data []int) chan int {
 out := make(chan int)
 go func() {
  for _, v := range data {
   out <- v
  }
  close(out)
 }()
 return out
}

// Recebe um canal de inteiros e retorna um canal de inteiros com os valores pares
func filter(in chan int) chan int {
 out := make(chan int)
 go func() {
  for v := range in {
   if v%2 == 0 {
    out <- v
   }
  }
  close(out)
 }()
 return out
}

// Recebe um canal de inteiros e retorna um canal de inteiros com os valores ao quadrado
func square(in chan int) chan int {
 out := make(chan int)
 go func() {
  for v := range in {
   out <- v * v
  }
  close(out)
 }()
 return out
}

func main() {
 data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

 c1 := generate(data)
 c2 := filter(c1)
 c3 := square(c2)

 for v := range c3 {
  fmt.Println(v)
 }
}

```

Saida

```bash
go run pipeline.go 
4
16
36
64
```

### Fan-In - Fan-Out

Este patterns permite que um mesmo canal possa ser compartilhado entre várias go routines (Fan-Out) que, diferente do pattern Worker Pools, trabalharam os dados de forma independente, podendo gerar dados distintos. Estes dados são agreggados ao final da execução (Fan-In).

```go
func cymbal(drum chan string) {
 drum <- "-tsssss-"
}

func hiHat(drum chan string) {
 drum <- "-chik-chik-"
}

func snare(drum chan string) {
 drum <- "-tak-"
}

func bassDrum(drum chan string) {
 drum <- "-boom-"
}

func main() {
 drumSet := make(chan string)
 var song string

 go hiHat(drumSet)
 go bassDrum(drumSet)
 go snare(drumSet)
 go cymbal(drumSet)

 for i := 0; i < 4; i++ {
  song += <-drumSet
 }
 fmt.Println(song)

}

```

Saida

```bash
go run fanin-fanout.go 
-tsssss--tak--chik-chik--boom-
```

### Menção Honrosa

Aqui vale fazer uma menção honrosa a alguns padrões fundamentais de concorrência em Go

- Generator - O padrão generator é, basicamente, uma função que roda uma go routine anônima e retorna um channel para comunicação com outras rotinas.

    ```go
    func writer(text string) <-chan string {
     c := make(chan string)
     go func() {
      for {
       c <- fmt.Sprintf("Hello %s", text)
      }
     }()
     return c
    }
    ```

- WaitGroup - Utilizando a entidade WaitGroup, do pacote sync, é possível determinar a quantidade de go routines que serão executadas e esperar cada uma delas finalizar sua tarefa para continuar a execução do código.

    ```go
    func writer(text string, wg *sync.WaitGroup) {
     fmt.Printf("Hello %s\n", text)
     wg.Done() // notifica a finalização da execução
    }
    func main() {
     var wg sync.WaitGroup
     wg.Add(3) // o wait group vai esperar a finalização de 3 goroutines
     go writer("world", &wg)
     go writer("golang", &wg)
     go writer("universe", &wg)
     wg.Wait()// aguarda a finalização das 3 rotinas
     fmt.Println("All done!")
    }
    ```

- Multiplex - Permite que uma goroutine receba os dados de vários channels ao mesmo tempo usando a instrução Select.

    ```go
    select {
    case msg1 := <-ch1:
        // processa a mensagem do channel ch1
    case msg2 := <-ch2:
        // processa a mensagem do channel ch2
    case <-time.After(10 * time.Second):
        // executa a ação de timeout
    }
    ```
