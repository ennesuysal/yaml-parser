package main

func (d *diagnostic) diagArrContStr(line string, indent float32) {
	d.continuingStr = continuingString{}
	d.continuingStrIndent = indent
	d.parseArrayCont(parseArrContStr(line)+":", d.continuingArrIndent)
	d.continuingStrRoot = d.root[len(d.root)-1][0].(*node)
	d.lastContString = nil
}
