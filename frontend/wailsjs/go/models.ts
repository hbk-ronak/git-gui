export namespace types {
	
	export class Branch {
	    Name: string;
	    IsCurrent: boolean;
	    IsRemote: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Branch(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.IsCurrent = source["IsCurrent"];
	        this.IsRemote = source["IsRemote"];
	    }
	}
	export class CommitResult {
	    Success: boolean;
	    CommitSHA: string;
	    Message: string;
	
	    static createFrom(source: any = {}) {
	        return new CommitResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Success = source["Success"];
	        this.CommitSHA = source["CommitSHA"];
	        this.Message = source["Message"];
	    }
	}
	export class DiffHunk {
	    Header: string;
	    OldStart: number;
	    OldLines: number;
	    NewStart: number;
	    NewLines: number;
	    Lines: string[];
	
	    static createFrom(source: any = {}) {
	        return new DiffHunk(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Header = source["Header"];
	        this.OldStart = source["OldStart"];
	        this.OldLines = source["OldLines"];
	        this.NewStart = source["NewStart"];
	        this.NewLines = source["NewLines"];
	        this.Lines = source["Lines"];
	    }
	}
	export class DiffResult {
	    FilePath: string;
	    Diff: string;
	    Hunks: DiffHunk[];
	
	    static createFrom(source: any = {}) {
	        return new DiffResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.FilePath = source["FilePath"];
	        this.Diff = source["Diff"];
	        this.Hunks = this.convertValues(source["Hunks"], DiffHunk);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class FileStatus {
	    Path: string;
	    Status: string;
	    Staged: boolean;
	
	    static createFrom(source: any = {}) {
	        return new FileStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Path = source["Path"];
	        this.Status = source["Status"];
	        this.Staged = source["Staged"];
	    }
	}
	export class GitRepo {
	    Path: string;
	    CurrentBranch: string;
	
	    static createFrom(source: any = {}) {
	        return new GitRepo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Path = source["Path"];
	        this.CurrentBranch = source["CurrentBranch"];
	    }
	}

}

