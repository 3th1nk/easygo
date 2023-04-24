package charset

import "github.com/3th1nk/easygo/charset/internal"

/*
	GB2312
	1980年，中国发布了第一个汉字编码标准，也即GB2312，全称《信息交换用汉字编码字符集·基本集》，通常简称GB（“国标”汉语拼音首字母），共收录了 6763 个常用的汉字和字符，此标准于次年5月实施，它满足了日常 99% 汉字的使用需求

	GBK
	由于有些汉字是在GB2312标准发布之后才简化的，还有一些人名、繁体字、日语和朝鲜语中的汉字也没有包括在内，所以，在GB2312的基础上添加了这部分字符，就形成了GBK，全称《汉字内码扩展规范》，共收录了两万多个汉字和字符，它完全兼容GB2312
	GBK于1995年发布，不过它只是 "技术规范指导性文件"，并不属于国家标准

	GB18030
	GB18030全称《信息技术 中文编码字符集》，共收录七万多个汉字和字符，它在GBK的基础上增加了中日韩语中的汉字和少数名族的文字及字符，完全兼容GB2312，基本兼容GBK
	GB18030发布过两个版本，第一版于2000年发布，称为GB18030-2000，第二版于2005年发布，称为 GB18030-2005

	简单概括：
		1.GB2312兼容ASCII, GBK兼容GB2312，GB18030兼容GB2312和GBK
		2.GB2312、GBK、GB18030以及UTF8共同点是都兼容ASCII

	GB2312把每个汉字都编码成两个字节，第一个字节是高位字节，第二个字节是低位字节，编码范围如下:
	字节数	  码位区间				 编码数
			 第一字节     第二字节
	单字节	  0x00-0x7F				128				ASCII的编码范围
	双字节    0xA1-0xA9  0xA1-0xFE  846				GB2312编码范围
			 0xB0-0xF7  0xA1-0xFE  6768			   GB2312编码范围

	GBK也是双字节编码，为了向下兼容GB2312，GBK使用了GB2312没有用到的编码区域，总的编码范围是: 第一个字节 0x81–0xFE，第二个字节 0x40–0xFE，具体的编码范围如下：
	字节数	  码位区间				 		  编码数
			 第一字节     第二字节
	单字节	  0x00-0x7F						 128		ASCII的编码范围
	双字节    0xA1-0xA9  0xA1-0xFE  		 846		GB2312编码范围
			 0xB0-0xF7  0xA1-0xFE  			6768	   GB2312编码范围
			 0x81-0xA0  0x40-0xFE(7F除外)	 6080	    GBK编码范围
			 0xA8-0xA9  0x40-0xA0(7F除外)    192		GBK编码范围
			 0xAA-0xFE  0x40-0xA0(7F除外)    8160		GBK编码范围
			 0xA1-0xA7  0x40-0xA0(7F除外)    自定义	     用户自定义编码范围
			 0xAA-0xAF  0xA1-0xFE			自定义		用户自定义编码范围
			 0xF8-0xFE  0xA1-0xFE			自定义		用户自定义编码范围


	GB18030是变长多字节字符集，每个字或字符可以由一个，两个或四个字节组成，所以它的编码空间是很大的，最多可以容纳161万个字符，由于需要兼容 GBK，四个字节的前两个字节和GBK编码保持一致，具体的编码范围如下：
	字节数	  码位区间				 		  				编码数
			 第一字节     第二字节	   第三字节	    第四字节
	单字节	  0x00-0x7F						 				128			ASCII的编码范围
	双字节    0x81-0xFE  0x40-0x7E	     	 				23940	    GBK编码范围
					    0x80-0xFE
	四字节	  0x81-0xFE  0x30-0x39	 0x81-0xFE  0x30-0x39  1587600 	    GB18030编码范围

*/

func IsGB2312(s string) bool {
	var nBytes int
	for _, chr := range s {
		// 跳过ASCII字符
		if nBytes == 0 && chr < 0x80 {
			continue
		}

		if nBytes == 0 {
			if chr < 0xA1 || chr > 0xF7 {
				return false
			}
			nBytes = 2
		} else {
			if chr < 0xA1 || chr > 0xFE {
				return false
			}
		}
		nBytes--
	}

	if nBytes != 0 {
		return false
	}
	return true
}

func IsGBK(s string) bool {
	var nBytes int
	for _, chr := range s {
		// 跳过ASCII字符
		if nBytes == 0 && chr < 0x80 {
			continue
		}

		if nBytes == 0 {
			if chr < 0x81 || chr > 0xFE {
				return false
			}
			nBytes = 2
		} else {
			if chr < 0x40 || chr > 0xFE || chr == 0x7F {
				return false
			}
		}
		nBytes--
	}

	if nBytes != 0 {
		return false
	}
	return true
}

func IsGB18030(s string) bool {
	var nBytes int
	for _, chr := range s {
		// 跳过ASCII字符
		if nBytes == 0 && chr < 0x80 {
			continue
		}

		if nBytes == 0 {
			if chr < 0x81 || chr > 0xFE {
				return false
			}
			nBytes = 2
		} else {
			if chr < 0x30 || chr > 0xFE || chr == 0x7F {
				return false
			}
		}
		nBytes--
	}

	if nBytes != 0 {
		return false
	}
	return true
}

func ToGBK(s string) string {
	if IsUTF8(s) {
		return internal.Translate(s, "UTF-8", "GBK")
	}
	return s
}

func ToGB18030(s string) string {
	if IsUTF8(s) {
		return internal.Translate(s, "UTF-8", "GB18030")
	}
	return s
}
