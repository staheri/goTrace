
const (
	_grtnNode           = 0 // Goroutine events: EvGoCreate, EvGoStart, EvGoEnd
	_chnlNode           = 1
)

type Node struct{
	id         int
	typ        int
	g          int
	eid        int
	label      string
	color      string
	posx       float64
	posy       float64
	bold       bool
	style      string
}

type Edge struct{
	id       int
	src      *Node
	dest     *Node
	color    string
	bold     bool
	style    string
}



func ToFile(dbName string){
	// Establish connection
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/"+dbName)
	if err != nil {
		fmt.Println(err)
	}else{
		fmt.Println("Connection Established")
	}
	defer db.Close()

	//goEvents   := []string{"EvGoCreate","EvGoStart","EvGoEnd"}
	//chanEvents := []string{"EvChMake","EvChSend","EvChRecv","EvChClose"}

	// Vars
	var q string
	var id,gid,parent_id,end_eid,create_eid,start_eid int

	// Store nodes
	//var nodes []*Nodes
	//data := make(map[int]*Node) // key: eventID



	// GoCreate/GoStart nodes
	q = `SELECT id, gid, parent_id, ended, create_eid, start_eid
	     FROM Goroutines;`
	fmt.Printf(">>> Executing %s...\n",q)
	res, err := db.Query(q)
	if err != nil {
		panic(err)
	}

	for res.Next(){
		err = res.Scan(&id, &gid, &parent_id, &end_eid, &create_eid, &start_eid)
		//fmt.Printf("Goroutine %d created by %d @ %d\n",child_g,parent_g,eid)
		if id == 1{
			// gid=0 then there is no creator
		} else{
			//
		}
	}


	// Query and generate nodes
  // >> Goroutine eventes (create, start, end, ...)
	     // count: select Count(*), t1.g from events t1 INNER JOIN global.catGRTN t3 ON t1.type = t3.eventName GROUP BY t1.g;
			 // event list: select t1.id, t1.type, t1.g from events t1 INNER JOIN global.catGRTN t3 ON t1.type = t3.eventName WHERE t1.type="EvGoCreate" OR t1.type="EvGoStart" OR t1.type="EvGoEnd";
	// >> Channel events (make, send, recv, close)
	// >> (Other events: locks, waitingGroups, GC, PROC, etc.)

	// Query and generate edges
	// >>

}
