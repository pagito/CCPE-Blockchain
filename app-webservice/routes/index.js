var express = require('express');
var router = express.Router();

var app = express();
var http = require('http').Server(app);

var my_cc;

 // Step 1 ================================== 
    var Ibc1 = require('ibm-blockchain-js');
    var ibc = new Ibc1(/*logger*/);             //you can pass a logger such as winston here - optional 
    var chaincode = {};
 
    // ================================== 
    // configure ibc-js sdk 
    // ================================== 
    var options =   {
        network:{
            peers:   [{
                "api_host": "172.17.0.2",
                "api_port": 7050,
                //"api_port_tls": xxx,
                //"id": "xxxxxx-xxxx-xxx-xxx-xxxxxxxxxxxx_vpx"
                "id": "app_peer"
            }],
            users:  [{
                "enrollId": "test_user0",
                "enrollSecret": "MS9qrN8hFjlE"
            }],
            options: {              //this is optional 
                quiet: true, 
                timeout: 60000,
                tls: false,
            }
        },
        chaincode:{
            zip_url: 'https://github.com/pagito/CCPE-Blockchain/raw/master/chaincode/ccpe_chaincode.zip',
            unzip_dir: '/',
            git_url: 'https://github.com/pagito/CCPE-Blockchain/chaincode',
            deployed_name: null
            //,deployed_name:'5413191f18c5cab35639e42515edbb47c12c2ce7306d107b7cc6e23b591a5a4c123c261fd0da1fdc2214047e0d168b4087b86b3c86d4d62a219b46b9a1abc48e'
        }
    };
    
    // Step 2 ================================== 
    ibc.load(options, cb_ready);
 
    // Step 3 ==================================
    function cb_ready(err, cc){                             //response has chaincode functions
        //app1.setup(ibc, cc);
        //app2.setup(ibc, cc);

    // Step 4 ==================================
        if(err != null){
            console.log('! looks like an error loading the chaincode or network, app will fail\n', err);
            if(!process.error) process.error = {type: 'load', msg: err.details};                //if it already exist, keep the last error
        } else {
            
            my_cc = cc;

            if(!cc.details.deployed_name || cc.details.deployed_name === ''){                //decide if I need to deploy or not
                //console.log("Ready, but do not deploy yet");
                
                cc.deploy('init', ['99'], {delay_ms: 30000}, function(e){                    //delay_ms is milliseconds to wait after deploy for conatiner to start, 50sec recommended
                    console.log("success deployed");
                    cb_deployed();
                });
            }
            else{
                //my_cc = cc;
                console.log('chaincode summary file indicates chaincode has been previously deployed');
             
                cb_deployed();
            }
        }
    }

    // Step 5 ==================================
    function cb_deployed(err){
        console.log('sdk has deployed code and waited');
        my_cc.invoke.Init();
        console.log('Init done');
        my_cc.query.read(['dummy_query']);
        console.log('Query done');
        http.listen(8088, function(){
          console.log('listening on *:3000');
          
        });
    }


/* GET home page. */
router.get('/', function(req, res, next) {
  res.render('index', { title: 'Express' });
});


/* Example of GET method, simple respond with a string */
router.get('/transaction', function(req, res) {
  var responseObject =  { message: "Get Transaction Page" };
  res.send(responseObject);
});



/* Get POST data coming from Exchange APP */
router.post('/getTransaction', function(req, res, next) {
	console.log('savedata called');
    console.log(req.body);
    
    var id = req.body.Transaction_id;
    var userA = req.body.User_A;
    var userB = req.body.User_B;
    var sellerA = req.body.Seller_A;
    var sellerB = req.body.Seller_B;    
    var pointA = req.body.Ex_points_A;
    var pointB = req.body.Ex_points_B;
    //console.log(req.body.User_A);

	var curret_date = new Date();
    var dateStr = curret_date.getFullYear()+''+(curret_date.getMonth()+1)+''+curret_date.getDate();
    var tmpID = sellerA+'-'+sellerB+'-'+dateStr+'-'+id;
    my_cc.invoke.init_transaction([tmpID,userA,userB,sellerA,sellerB,pointA,pointB,''+Date.parse(new Date())],function(err, data) {

        var succ_data = data;
        res.json({
            "msg":succ_data,
            "respond":true,
            "record_id":id
        });
        console.log('success',succ_data);
    });

    var responseObject =  { "message": "Transaction accepted",
							"respond": true };
	res.send(responseObject);
});


module.exports = router;
