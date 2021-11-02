package visualizer

import (
	"log"
	"runtime"
	"strconv"

	"../algorithms"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/nullboundary/glfont"
)

const (
	width  = 1280
	height = 700
	rows   = 200
)

var (
	columns   = 0
	rectangle = []float32{
		-0.1, 0.1, 0,
		-0.1, -0.1, 0,
		0.1, -0.1, 0,

		-0.1, 0.1, 0,
		0.1, 0.1, 0,
		0.1, -0.1, 0,
	}
	window *glfw.Window
	font   *glfont.Font
)

type bar struct {
	drawable uint32
}

type graph struct {
	bars      []*bar
	color     []bool
	yPosition float32
	side      bool
	done      bool
}

//generateList: Crea una lista de numeros aleatorios con el metodo de
//congruncia lineal multiplicativa
//(x): Semilla. Debe de ser primo entre [11, 101]  0 <= x < m
//(n): Cantidad de numeros aleatorios.
//(m): Periodo
func generateList(n int, x int, m int) []int {
	var a int = 1103515245 //    0 < a < m multiplicador
	var c int = 12345      //      0 <= c < m  Incremento
	var nums []int
	for i := 0; i < n; i++ {
		x = (a*x + c) % m
		nums = append(nums, x%31)
	}
	return nums
}

//Start ...
func Start(n int, x int, m int, msgCh chan string) {
	//Constantes
	algorithmsName := [6]string{"BubbleSort", "SelectionSort", "InsertionSort", "MergeSort", "QuickSort", "HeapSort"}
	stopCh := make(chan struct{}) //Canal para detener todo

	//GENERA LA LISTA DE NUMEROS ALEATORIOS
	numberList := generateList(n, x, m)
	columns = len(numberList) + int(float32(len(numberList))*0.05)

	var channelList []chan [][]int //Lista de canales
	var graphList []*graph         //Lista de graficos
	var numberLists [][]int        //Lista de listas de numeros

	//INICIALIZA TODOS LOS DATOS
	for i := 0; i < 6; i++ {
		newList := make([]int, len(numberList), len(numberList)) //Crea una lista nueva
		copy(newList, numberList)                                //Copia los datos de la lista de aleatorios en la lista nueva
		numberLists = append(numberLists, newList)               //Guarda la lista en una lista de listas

		channelList = append(channelList, make(chan [][]int)) //Crea un canal nuevo y lo guarda en la lista de canales
	}

	//INICIA CADA ALGORITMO CON CORRUTINAS
	msgCh <- "Inicia corrida con " + strconv.Itoa(len(numberList)) + " elementos aleatorios"
	go algorithms.HeapSort(numberLists[0], channelList[0], stopCh, msgCh)
	go algorithms.QuickSort(numberLists[1], channelList[1], stopCh, msgCh)
	go algorithms.MergeSort(numberLists[2], channelList[2], stopCh, msgCh)
	go algorithms.InsertionSort(numberLists[3], channelList[3], stopCh, msgCh)
	go algorithms.SelectionSort(numberLists[4], channelList[4], stopCh, msgCh)
	go algorithms.BubbleSort(numberLists[5], channelList[5], stopCh, msgCh)

	//INICIA LA VENTANA
	runtime.LockOSThread()
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	window = initGlfw()
	initOpenGL()

	//DIBUJA LAS ETIQUETAS Y CREA LOS GRAFICOS
	font, _ = glfont.LoadFont("Roboto-Light.ttf", int32(52), width, height)
	for i := 0; i < 6; i++ {
		//Los primero 3 graficos se ubican a la derecha, los otros 3 a la izquierda
		//Las primeras 3 eriquetas de ubican a la ziquierda, las otras 3 a la derecha
		x := 100
		side := false
		if i > 2 {
			x = 800
			side = true
		}
		//PINTA LAS ETIQUETAS EN AMBOS BUFFERS
		font.Printf(float32(x), (float32(i)+0.7)*120, 1.2, algorithmsName[i])
		window.SwapBuffers()
		font.Printf(float32(x), (float32(i)+0.7)*120, 1.2, algorithmsName[i])

		newGraph := createGraph(3.4*float32(i), numberLists[i], side) //Crea un grafico
		graphList = append(graphList, newGraph)                       //Guarda el grafico en la lista de grafico
	}

	gl.Enable(gl.SCISSOR_TEST)  //Activa la opcion de cortar la pantalla
	for !window.ShouldClose() { //Mientras la ventana no se cierre
		//Se recibe data(cambios) de los algoritmos para poder actualizar la pantalla
		select {
		case changesHP := <-channelList[0]: //HeapSort
			drawInWindow(0, 0, changesHP, graphList[0])
		case changesQS := <-channelList[1]: //QuickSort
			drawInWindow(0, 117, changesQS, graphList[1])
		case changesMS := <-channelList[2]: //MergeSort
			drawInWindow(0, 234, changesMS, graphList[2])
		case changesIS := <-channelList[3]: //InsertionSort
			drawInWindow(640, 351, changesIS, graphList[3])
		case changesSS := <-channelList[4]: //SelectionSort
			drawInWindow(640, 468, changesSS, graphList[4])
		case changesBS := <-channelList[5]: //BubbleSort
			drawInWindow(640, 585, changesBS, graphList[5])
		}
		glfw.PollEvents()
	}
	close(stopCh) //Cerrando este canal cerramos los demas canales en cada algoritmo
	close(msgCh)  //Cerramos el canal de mensajes
}

//drawInWindow: permite dibujar en pantalla
//(xCut,yCut): coordenadas iniciales para cortar.
//(changes): cambios que hizo el algoritmo.
//(graph): grafico correspondiente al algoritmo.
func drawInWindow(xCut, yCut int32, changes [][]int, graph *graph) {
	if !graph.done { //Si el grafico no se ha terminado de pintar
		gl.Scissor(xCut, yCut, 640, 117)                    //Seleccionamos la parte de la ventana que queremos actualizar
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT) //La limpiamos
		if len(changes) == 0 {                              //Si el algoritmo termino
			graph.done = true                      //El grafico se pinto
			graph.color = []bool{true, true, true} //Actualiza el color del grafico
		} else { //Si el algoritmo no ha terminado
			graph.updateGraph(changes) //Actualiza las barras
		}
		//Pinta el grafico en ambos buffers
		graph.drawGraph()
		window.SwapBuffers()
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT) //La limpiamos
		graph.drawGraph()
	}
}

//======================FUNCIONES DEl GRAFICO======================================

//createGraph: crea un grafico nuevo con todos sus atributos.
//(yPos): Posicion del grafico en el eje y
//(data): Lista de numeros que seran representados por el grafico
//(side): Lado en el cual se va a posicionar el algoritmo(derecha/true , izquierda/false)
func createGraph(yPos float32, data []int, side bool) *graph {
	var newBars []*bar
	for i := 0; i < len(data); i++ {
		newB := createBar(float32(i), yPos, data[i], side)
		newBars = append(newBars, newB)
	}

	newGraph := &graph{
		yPosition: yPos,
		color:     []bool{true, false, true},
		bars:      newBars,
		side:      side,
		done:      false,
	}
	return newGraph
}

//drawGraph: Recorre la lista de barras del grafico y pinta cada una.
func (g *graph) drawGraph() {
	for bar := 0; bar < len(g.bars); bar++ {
		g.bars[bar].drawBar(g.color)
	}
}

//updateGraph: Actualiza las barras
//(changes): Lista de cambios que realizo un algoritmo
func (g *graph) updateGraph(changes [][]int) {
	for change := 0; change < len(changes); change++ {
		index := changes[change][1] //En pos 1 del cambio esta el indice de la barra que quiero cambiar
		value := changes[change][0] //En pos 0 del cambio esta el valor nuevo que corresponde a ese indice

		//Le asigno a la barra en ese indice un nuevo drawable con los valores del cambio
		g.bars[index].setDrawable(float32(index), g.yPosition, value, g.side)
	}
}

//======================FUNCIONES DE LA BARRA======================================

//createBar: Crea una barra con su atributo dibujable
//(x,y): Posicion en el eje "x" y "y" de la barra
//(value): Valor correspondiente a esa barra
//(side): Lado donde se ubica esa barra
func createBar(x, y float32, value int, side bool) *bar {
	bar := bar{}
	bar.setDrawable(x, y, value, side)
	return &bar
}

//setDrawable: Genera y establece el drawable para la barra, es el objeto que se pintara.
//(x,y): Posicion en el eje "x" y "y" de la barra
//(value): Valor correspondiente a esa barra
//(side): Lado donde se ubica esa barra
func (c *bar) setDrawable(x, y float32, value int, side bool) {
	points := make([]float32, len(rectangle), len(rectangle))
	copy(points, rectangle)

	for i := 0; i < len(points); i++ {
		var position, size, m float32
		switch i % 3 {
		case 0: //Se trata de una x
			size = (1 / float32(columns))
			position = x * size / 2 // POSITION
			m = 1
			if side == false {
				m = 0
			}
		case 1: //Se trata de una y
			size = (float32(value) / float32(rows)) / 2
			position = 0
			m = y / 10
		default: //Se trata de una z
			continue
		}

		if points[i] < 0 {
			points[i] = ((position * 2) - 1) + m
		} else {
			points[i] = (((position + size) * 2) - 1) + m
		}
	}
	c.drawable = makeVao(points)
}

//drawBar: Dibuja la barra en la pantalla
//(color): color en que sera pintada la barra
func (c *bar) drawBar(color []bool) {
	gl.ColorMask(color[0], color[1], color[2], false)
	gl.BindVertexArray(c.drawable)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(rectangle)/3))
}

//Genera el objeto VAO con la liberia OpenGL
func makeVao(points []float32) uint32 {
	var vbo, vao uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	return vao
}

//======================FUNCIONES DE LA VENTANA======================================

func initGlfw() *glfw.Window {
	window, _ := glfw.CreateWindow(int(width), int(height), "THE BEST SORT VISUALIZER", nil, nil)
	window.MakeContextCurrent()
	return window
}
func initOpenGL() {
	if err := gl.Init(); err != nil {
		panic(err)
	}
}
