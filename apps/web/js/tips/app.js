/*公用弹窗 2019年12月27日 15:22:36 */
function add_load(){
	if($("div").hasClass("upload_loading")){
		$(".upload_loading").css("opacity",1);
		$(".upload_loading").show();
	}else{
		let html = '<style>.upload_loading{height:50px;width:50px;text-align:center;overflow:hidden;padding:10px;position:fixed;top:50%;left:50%;margin-top:-25px;margin-left:-25px;z-index:1000}.container_upload{width:30px;height:30px;position:relative}.container_upload.animation-upload{-webkit-transform:rotate(10deg);transform:rotate(10deg);-webkit-animation:rotation 1s infinite;animation:rotation 1s infinite}.container_upload.animation-upload .shape{border-radius:5px}.container_upload .shape{position:absolute;width:10px;height:10px;border-radius:1px}.container_upload .shape.shape1{left:0;background-color:#5c6bc0}.container_upload .shape.shape2{right:0;background-color:#8bc34a}.container_upload .shape.shape3{bottom:0;background-color:#ffb74d}.container_upload .shape.shape4{bottom:0;right:0;background-color:#f44336}@-webkit-keyframes rotation{from{-webkit-transform:rotate(0deg);transform:rotate(0deg)}to{-webkit-transform:rotate(360deg);transform:rotate(360deg)}}@keyframes rotation{from{-webkit-transform:rotate(0deg);transform:rotate(0deg)}to{-webkit-transform:rotate(360deg);transform:rotate(360deg)}}.animation-upload .shape1{-webkit-animation:animationshape1 .5s ease 0s infinite alternate;animation:animationshape1 .5s ease 0s infinite alternate}@-webkit-keyframes animationshape1{from{-webkit-transform:translate(0,0);transform:translate(0,0)}to{-webkit-transform:translate(20px,20px);transform:translate(20px,20px)}}@keyframes animationshape1{from{-webkit-transform:translate(0,0);transform:translate(0,0)}to{-webkit-transform:translate(20px,20px);transform:translate(20px,20px)}}.animation-upload .shape2{-webkit-animation:animationshape2 .5s ease 0s infinite alternate;animation:animationshape2 .5s ease 0s infinite alternate}@-webkit-keyframes animationshape2{from{-webkit-transform:translate(0,0);transform:translate(0,0)}to{-webkit-transform:translate(-20px,20px);transform:translate(-20px,20px)}}@keyframes animationshape2{from{-webkit-transform:translate(0,0);transform:translate(0,0)}to{-webkit-transform:translate(-20px,20px);transform:translate(-20px,20px)}}.animation-upload .shape3{-webkit-animation:animationshape3 .5s ease 0s infinite alternate;animation:animationshape3 .5s ease 0s infinite alternate}@-webkit-keyframes animationshape3{from{-webkit-transform:translate(0,0);transform:translate(0,0)}to{-webkit-transform:translate(20px,-20px);transform:translate(20px,-20px)}}@keyframes animationshape3{from{-webkit-transform:translate(0,0);transform:translate(0,0)}to{-webkit-transform:translate(20px,-20px);transform:translate(20px,-20px)}}.animation-upload .shape4{-webkit-animation:animationshape4 .5s ease 0s infinite alternate;animation:animationshape4 .5s ease 0s infinite alternate}@-webkit-keyframes animationshape4{from{-webkit-transform:translate(0,0);transform:translate(0,0)}to{-webkit-transform:translate(-20px,-20px);transform:translate(-20px,-20px)}}@keyframes animationshape4{from{-webkit-transform:translate(0,0);transform:translate(0,0)}to{-webkit-transform:translate(-20px,-20px);transform:translate(-20px,-20px)}}</style><div class="upload_loading"><div class="container_upload animation-upload"><div class="shape shape1"></div><div class="shape shape2"></div><div class="shape shape3"></div><div class="shape shape4"></div></div></div>';
		$(document.body).append(html);
		$(".upload_loading").show();
	}
}
function remove_load(){
	$(".upload_loading").animate({"opacity":0}, 500);
	setTimeout(function(){$(".upload_loading").hide();},500);
}