package workflow

import (
	"reflect"
	"testing"
)

func TestUnimplementedITask_DryRun(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "dry run",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UnimplementedITask{}
			if err := u.DryRun(); (err != nil) != tt.wantErr {
				t.Errorf("DryRun() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUnimplementedITask_GetWorkflow(t *testing.T) {
	tests := []struct {
		name string
		want *Workflow
	}{
		{
			name: "get workflow",
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UnimplementedITask{}
			if got := u.GetWorkflow(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetWorkflow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnimplementedITask_Start(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "start",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UnimplementedITask{}
			if err := u.Start(); (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUnimplementedITask_Stop(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "stop",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UnimplementedITask{}
			if err := u.Stop(); (err != nil) != tt.wantErr {
				t.Errorf("Stop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUnimplementedITask_Pause(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "pause",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UnimplementedITask{}
			if err := u.Pause(); (err != nil) != tt.wantErr {
				t.Errorf("Pause() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUnimplementedITask_Resume(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "resume",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UnimplementedITask{}
			if err := u.Resume(); (err != nil) != tt.wantErr {
				t.Errorf("Resume() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUnimplementedITask_SetParams(t *testing.T) {
	type args struct {
		params *TaskParams
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "set params",
			args: args{
				params: &TaskParams{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UnimplementedITask{}
			if err := u.SetParams(tt.args.params); (err != nil) != tt.wantErr {
				t.Errorf("SetParams() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
