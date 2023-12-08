package util_test

import (
	"github.com/acgtools/hanime-hunter/pkg/util"
	"golang.org/x/net/html"
	"strings"
	"testing"
)

const doc = `

<!DOCTYPE HTML>
<html lang="en">

<head>
	<meta charset="utf-8">
	<meta http-equiv="X-UA-Compatible" content="IE=edge">
	<meta name="referrer" content="no-referrer">
	<meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1">
	<meta name="google-site-verification" content="oV77KljbCxlciy-aD-Uy_dZSYUENVR_6jAhWSp_cb48" />
	<meta name="exoclick-site-verification" content="c02975e2897725fd5f30045bf364309a">
	<meta name="juicyads-site-verification" content="cc330848f3dfc20e8259699c6a096411">
	<link rel="shortcut icon" type="image/x-icon" href="https://img4.qy0.ru/data/2205/36/tab_logo.png" />
	<link rel="apple-touch-icon" href="https://img4.qy0.ru/data/2205/36/tab_logo.png" />
	<link rel="canonical" href="https://hanime1.me/download?v=84720" />
	<meta name="RATING" content="RTA-5042-1996-1400-1577-RTA" />

	<meta property="og:url" content="https://hanime1.me/download?v=84720" />
	<meta property="og:type" content="article" />
	<meta property="og:title" content="下載 [中字後補] 魔騎夜談 2 - Hanime1.me" />
	<meta property="og:description" content="ナイトテール 第二話 この人、夢に出てきた人だ！ [中字後補] 夢魔を倒すもの、ナイトテールとして人々を守っている夜。
ちょっぴり恥ずかしい目に遭ったものの、今回も無事に夢魔を倒すことができたのだった。
夢魔退治で疲れた夜は、火照った身体を鎮めるべく、自らの身体を慰める。
はしたないとは思いつつも好奇心を抑えきれないお年頃。
その指先は激しくなる一方だった。
そんなことをしていればエッチな夢も見ちゃうもの。
夜の旺盛な好奇心は、夢でお尻の穴でのプレイを再現する。
当たり前の羞恥心、止めどない性への興味、そして感じたことのない背徳感。
彼女はますます深みにハマってしまう。
ところで、この人前にどこかで見なかった？
「あら、改めて言われると恥ずかしいですね」
ドキドキな夢魔退治は「夜」にお任せ！
でも……「夜」の中にいるモンモンな淫魔はどうすればいいの？" />
	<meta property="og:image"
		content="https://vdownload.hembed.com/image/thumbnail/WdQ34oOh.jpg?secure=JV7-DzeYCpjCVWPL8HfiAA==,1703989809" />

	<title>下載 [中字後補] 魔騎夜談 2&nbsp;-&nbsp;H動漫/裏番/線上看&nbsp;-&nbsp;Hanime1.me</title>
	<meta name="title" content="下載 [中字後補] 魔騎夜談 2 - H動漫/裏番/線上看 - Hanime1.me">
	<meta name="description" content="ナイトテール 第二話 この人、夢に出てきた人だ！ [中字後補] 夢魔を倒すもの、ナイトテールとして人々を守っている夜。
ちょっぴり恥ずかしい目に遭ったものの、今回も無事に夢魔を倒すことができたのだった。
夢魔退治で疲れた夜は、火照った身体を鎮めるべく、自らの身体を慰める。
はしたないとは思いつつも好奇心を抑えきれないお年頃。
その指先は激しくなる一方だった。
そんなことをしていればエッチな夢も見ちゃうもの。
夜の旺盛な好奇心は、夢でお尻の穴でのプレイを再現する。
当たり前の羞恥心、止めどない性への興味、そして感じたことのない背徳感。
彼女はますます深みにハマってしまう。
ところで、この人前にどこかで見なかった？
「あら、改めて言われると恥ずかしいですね」
ドキドキな夢魔退治は「夜」にお任せ！
でも……「夜」の中にいるモンモンな淫魔はどうすればいいの？">

	<meta name="title" content="Hanime1.me - H動漫/裏番/線上看">
	<title>Hanime1.me - H動漫/裏番/線上看</title>
	<meta name="description" content="Hanime1.me 帶給你最完美的H動漫、H動畫、裏番、里番、成人色情卡通片的線上看體驗，絕對沒有天殺的片頭廣告！">

	<!-- CSRF Token -->
	<meta name="csrf-token" content="mT40Kc9KtzqxlIUZaIq7CkvnVifgwQIqi24TRqun">

	<!-- Styles -->
	<link href="/css/app.css?id=bccf8b8fa56b630fefde" rel="stylesheet">
	<link href="https://fonts.googleapis.com/icon?family=Material+Icons|Material+Icons+Outlined|Material+Icons+Sharp"
		rel="stylesheet">
	<link href='https://fonts.googleapis.com/css2?family=Encode+Sans+Condensed:wght@700&display=swap' rel='stylesheet'>
	<link href="https://cdnjs.cloudflare.com/ajax/libs/OwlCarousel2/2.3.4/assets/owl.carousel.min.css" rel="stylesheet">
	<link href="https://cdnjs.cloudflare.com/ajax/libs/OwlCarousel2/2.3.4/assets/owl.theme.default.min.css"
		rel="stylesheet">

	<!-- Scripts -->
	<script src="https://cdn.jsdelivr.net/npm/js-cookie@3.0.0/dist/js.cookie.min.js"></script>
	<script src="https://code.jquery.com/jquery-3.3.1.min.js"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/OwlCarousel2/2.3.4/owl.carousel.min.js"></script>
</head>

<body>
	<div style="overflow-x: hidden; ">
		<div id="main-nav" style="z-index: 10000 !important; " class="main-nav-video-show hidden-xs">
			<a href="/" style="padding-right: 2.5%; color: white; font-size: 1.4em;">
				<span style="color: crimson">H</span>anime1<span style="color: crimson">.</span>me
			</a>
			<a class="nav-item hidden-xs nav-desktop-items " href="https://hanime1.me/search?genre=裏番">裏番</a>
			<a class="nav-item hidden-xs hidden-sm nav-desktop-items " href="/previews/202312">新番預告</a>
			<a class="nav-item hidden-xs nav-desktop-items " href="https://hanime1.me/search?genre=泡麵番">泡麵番</a>
			<a class="nav-item hidden-xs nav-desktop-items " href="https://hanime1.me/search?genre=Motion+Anime">Motion
				Anime</a>
			<a class="nav-item hidden-xs nav-desktop-items " href="https://hanime1.me/search?genre=3D動畫">3D動畫</a>
			<a class="nav-item hidden-xs nav-desktop-items " href="https://hanime1.me/search?genre=同人作品">同人作品</a>
			<a class="nav-item hidden-xs nav-desktop-items nav-desktop-cosplay-item "
				href="https://hanime1.me/search?genre=Cosplay">Cosplay</a>
			<a class="nav-item hidden-xs hidden-sm hidden-md nav-desktop-items" href="https://hanime1.me/comics">H漫畫</a>
			<a class="nav-item hidden-xs hidden-sm nav-desktop-items" href="https://l.erodatalabs.com/s/0ZIRw4"
				target="_blank">無碼黃油</a>
			<!-- <a class="nav-item hidden-xs hidden-sm" href="https://hanime1.me/playlists">我的清單</a> -->

			<a style="padding-right: 0px; padding-left: 10px" class="nav-icon pull-right"
				href="https://hanime1.me/login"><span style="vertical-align: middle; font-size: 28px" class="material-icons-outlined">account_circle</span></a>

			<a class="nav-icon pull-right no-select"
				href="https://hanime1.me/search"><span style="vertical-align: middle; font-size: 28px;" class="material-icons-outlined">search</span></a>

			<a class="nav-icon pull-right no-select"
				href="/previews/202312"><span style="vertical-align: middle; font-size: 25px" class="material-icons-outlined">cast</span></a>
		</div>

		<div id="main-nav-home"
			style="z-index: 10001; padding:0; height: 52px; line-height: 52px; position: fixed; background-color: black"
			class="hidden-sm hidden-md hidden-lg nav-main-mobile">
			<div id="main-nav-home-mobile"
				style="z-index: 10000 !important; position: fixed !important; overflow-x: hidden; background: none; transition: height 0.3s, background-color 0.4s, backdrop-filter 0.4s, -webkit-backdrop-filter 0.4s, top 0.4s; height: 52px !important; overflow-y: hidden; overflow-x: hidden;"
				class="hidden-sm hidden-md hidden-lg">

				<div style="padding: 0 15px;">
					<a href="/"
						style="padding-right: 2.5%; color: white; font-size: 1.40em; line-height: 57px; margin-left: 5px;">
						<img style="width: 15px; margin-top: -7px; margin-right: 2px;" src="https://vdownload.hembed.com/image/icon/nav_logo.png?secure=HxkFdqiVxMMXXjau9riwGg==,4855471889">
    </a>

						<a id="user-modal-trigger" href="https://hanime1.me/login"
							style="padding-left: 1px; padding-right: 0px; cursor: pointer;"
							class="nav-icon pull-right no-select">
							<img style="width: 24px; border-radius: 2px;" src="https://vdownload.hembed.com/image/icon/user_default_image.jpg?secure=ue9M119kdZxHcZqDPrunLQ==,4855471320">
      </a>

							<a style="margin-top: -1px; padding: 0 11px;" class="nav-icon pull-right"
								href="https://hanime1.me/search"><img style="width: 31px;" src="https://vdownload.hembed.com/image/icon/search_input_placeholder.png?secure=10N-U1uEz-5YMgWwuLCfPw==,4855472065"></a>

								<a style="padding: 0 10px;" class="nav-icon pull-right"
									href="/previews/202312"><span style="vertical-align: middle; font-size: 24px" class="material-icons-outlined">cast</span></a>
				</div>
			</div>
		</div>
		<div style="overflow-x: hidden;">

			<div id="content-div">
				<div class="row no-gutter video-show-width download-panel">

					<div class="hidden-xs hidden-sm" style="padding-bottom: 15px; text-align: center;">
						<!-- JuicyAds v3.1 -->
						<script type="text/javascript" data-cfasync="false" async
							src="https://poweredby.jads.co/js/jads.js"></script>
						<ins id="906694" data-width="908" data-height="270"></ins>
						<script type="text/javascript" data-cfasync="false" async>
							(adsbyjuicy = window.adsbyjuicy || []).push({'adzone':906694});
						</script>
						<!--JuicyAds END-->
					</div>

					<div class="hidden-xs hidden-md hidden-lg"
						style="text-align: center; padding-top: 10px; padding-bottom: 10px;">
						<span style="vertical-align: top; margin-top: 5px;">
		<!-- JuicyAds v3.1 -->
		<script type="text/javascript" data-cfasync="false" async src="https://poweredby.jads.co/js/jads.js"></script>
		<ins id="940482" data-width="300" data-height="262"></ins>
		<script type="text/javascript" data-cfasync="false" async>(adsbyjuicy = window.adsbyjuicy || []).push({'adzone':940482});</script>
		<!--JuicyAds END-->
	</span>

						<span style="vertical-align: top; margin-top: 5px;">
		<!-- JuicyAds v3.1 -->
		<script type="text/javascript" data-cfasync="false" async src="https://poweredby.jads.co/js/jads.js"></script>
		<ins id="940482" data-width="300" data-height="262"></ins>
		<script type="text/javascript" data-cfasync="false" async>(adsbyjuicy = window.adsbyjuicy || []).push({'adzone':940482});</script>
		<!--JuicyAds END-->
	</span>
					</div>

					<div class="hidden-sm hidden-md hidden-lg" style="text-align: center;">
						<!-- JuicyAds v3.1 -->
						<script type="text/javascript" data-cfasync="false" async
							src="https://poweredby.jads.co/js/jads.js"></script>
						<ins id="906695" data-width="300" data-height="112"></ins>
						<script type="text/javascript" data-cfasync="false" async>
							(adsbyjuicy = window.adsbyjuicy || []).push({'adzone':906695});
						</script>
						<!--JuicyAds END-->
					</div>
					<div class="col-md-12" style="background-color: #141414;">

						<div style="color: white" class="mobile-padding">
							<div>
								<div style="margin-bottom: -6px;">
									<p style="font-size: 12px; color: #bdbdbd; font-weight: 500">2023-11-26
										<span style="font-weight: normal;">&nbsp;|&nbsp;</span> 45.6萬次點閱</p>
								</div>
								<h3
									style="line-height: 30px; font-weight: bold; font-size: 1.5em; margin-top: 0px; color: white; margin-bottom: 15px;">
									[中字後補] 魔騎夜談 2</h3>
								<img class="download-image" src="https://vdownload.hembed.com/image/thumbnail/WdQ34oOh.jpg?secure=JV7-DzeYCpjCVWPL8HfiAA==,1703989809">

								<div class="hidden-sm hidden-md hidden-lg"
									style="text-align: center; padding-bottom: 15px;">
									<ins class="adsbyexoclick" data-zoneid="4372480"></ins>
								</div>

								<table class="download-table">
									<tr>
										<th></th>
										<th>影片畫質</th>
										<th>檔案類型</th>
										<th class="hidden-xs">檔案大小</th>
										<th>下載鏈結</th>
									</tr>
									<tr>
										<td style="text-align: center;">
											<span style="vertical-align: middle;" class="material-icons">play_circle_filled</span>
										</td>
										<td>
											標準畫質 (480p)
										</td>
										<td>mp4</td>
										<td class="hidden-xs">N/A</td>
										<td><a class="exoclick-popunder"
												style="text-decoration: none; color: white; text-align: center; background-color: crimson; padding: 5px 10px; border-radius: 5px;"
												href="https://vdownload-45.sb-cd.com/1/5/15197064-480p.mp4?secure=beYk_hjhQ0h_W1V4I2p_4A,1702020698&amp;m=45&amp;d=1&amp;_tid=15197064"
												download="[中字後補] 魔騎夜談 2">下載</a></td>
									</tr>
									<tr>
										<td style="text-align: center;">
											<span style="vertical-align: middle;" class="material-icons">play_circle_filled</span>
										</td>
										<td>
											低清畫質 (240p)
										</td>
										<td>mp4</td>
										<td class="hidden-xs">N/A</td>
										<td><a class="exoclick-popunder"
												style="text-decoration: none; color: white; text-align: center; background-color: crimson; padding: 5px 10px; border-radius: 5px;"
												href="https://vdownload-45.sb-cd.com/1/5/15197064-240p.mp4?secure=beYk_hjhQ0h_W1V4I2p_4A,1702020698&amp;m=45&amp;d=1&amp;_tid=15197064"
												download="[中字後補] 魔騎夜談 2">下載</a></td>
									</tr>
								</table>
							</div>
						</div>
					</div>

					<script type="application/javascript">
						(function() {

			    //version 1.0.0

			    var adConfig = {
			    "ads_host": "a.pemsrv.com",
			    "syndication_host": "s.pemsrv.com",
			    "idzone": 4213656,
			    "popup_fallback": false,
			    "popup_force": false,
			    "chrome_enabled": true,
			    "new_tab": false,
			    "frequency_period": 5,
			    "frequency_count": 1,
			    "trigger_method": 2,
			    "trigger_class": "exoclick-popunder",
			    "trigger_delay": 0,
			    "only_inline": false
			};

</body>

</html>

`

func BenchmarkFindTagByNameAttrs(b *testing.B) {
	d, _ := html.Parse(strings.NewReader(doc))

	for i := 0; i < b.N; i++ {
		util.FindTagByNameAttrs(d, "table", false, nil)
	}
}

func BenchmarkFindTagByRegExp(b *testing.B) {
	const re = "(<table[^>]*>)(.*?)(.*)(</table>)"

	for i := 0; i < b.N; i++ {
		util.FindTagByRegExp(doc, re)
	}
}
